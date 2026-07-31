package eval

import (
	"llewellyn-kevin/yak/ast"
)

func Eval(p *ast.Program) (*EvalState, []error) {
	state := NewEvalState()
	state.AddFunctionsFromMap(p.FunctionTable)
	errors := state.executeBlock(p.MainBlock)
	return state, errors
}

func PartialEval(p *ast.Program, state *EvalState) []error {
	if state == nil {
		state = NewEvalState()
	}
	state.AddFunctionsFromMap(p.FunctionTable)
	errors := state.executeBlock(p.MainBlock)
	return errors
}

func (e *EvalState) executeBlock(block *ast.Block) (errorList []error) {
	for _, expr := range block.Expressions {
		switch typed := expr.(type) {
		case ast.ExecuteBlockExpression:
			errorList = append(errorList, e.execBlockExpr(typed, block)...)
		case ast.ExecuteConditionalExpression:
			errorList = append(errorList, e.execCondExpr(typed, block)...)
		default:
			if err := e.evalExpression(expr); err != nil {
				errorList = append(errorList, err...)
			}
		}
	}
	return
}

func (e *EvalState) execBlockExpr(expr ast.ExecuteBlockExpression, block *ast.Block) []error {
	nested, err := block.GetBlock(expr.BlockId)
	if err != nil {
		return []error{err}
	}
	return e.executeBlock(nested)
}

func (e *EvalState) execCondExpr(expr ast.ExecuteConditionalExpression, block *ast.Block) []error {
	cond, err := block.GetConditional(expr.ConditionalId)
	if err != nil {
		return []error{err}
	}
	stack, err := e.ActiveStack()
	if err != nil {
		return []error{err}
	}
	val, err := stack.Pop()
	if err != nil {
		return []error{RuntimeError{Message: "tried to run a conditional branch on an empty stack"}}
	}
	if val.IsTruthy() {
		return e.executeBlock(cond.WhenTrue)
	}
	return e.executeBlock(cond.WhenFalse)
}

func (e *EvalState) evalExpression(expr ast.Expression) (errors []error) {
	var err error
	switch {
	case ast.IsLiteral(expr):
		err = e.evalLiteral(expr)
	case ast.IsUnaryOperator(expr):
		err = e.evalUnaryOperator(expr)
	case ast.IsBinaryOperator(expr):
		err = e.evalBinaryOperator(expr)
	case ast.IsBitwiseOperator(expr):
		err = e.evalBitwiseOperator(expr)
	case ast.IsAssignment(expr):
		err = e.evalAssignment(expr)
	case ast.IsStackOperationExpression(expr):
		err = e.evalStackOperation(expr)
	case ast.IsIdentifier(expr):
		errors = e.evalIdentifier(expr.(ast.IdentifierExpression))
	default:
		err = RuntimeErrorf("unknown expression %s", expr.String())
	}
	if err != nil {
		errors = append(errors, err)
	}
	return
}

func (e *EvalState) evalLiteral(expr ast.Expression) (err error) {
	stack, err := e.ActiveStack()
	if err != nil {
		return err
	}

	switch v := expr.(type) {
	case ast.IntLiteral:
		stack.Push(IntValue{Value: v.Value})
	case ast.FloatLiteral:
		stack.Push(FloatValue{Value: v.Value})
	case ast.StringLiteral:
		stack.Push(StringValue{Value: v.Value})
	case ast.SymbolLiteral:
		stack.Push(SymbolValue{Value: v.Value})
	case ast.BoolLiteral:
		stack.Push(BooleanValue{Value: v.Value})
	default:
		err = RuntimeErrorf("unknown literal type %s", expr.String())
	}

	return
}

func (e *EvalState) evalUnaryOperator(expr ast.Expression) (err error) {
	activeStack, err := e.ActiveStack()
	if err != nil {
		return
	}

	literal, err := activeStack.Pop()
	if err != nil {
		err = RuntimeErrorf("tried to use a unary operator on empty stack: '%s'", e.stackLabel())
		return
	}

	result, err := literal.DoUnaryOperation(expr)
	activeStack.Push(result)
	return
}

func (e *EvalState) evalStackOperation(expr ast.Expression) (err error) {
	stack, err := e.ActiveStack()
	if err != nil {
		return err
	}

	switch expr.(type) {
	case ast.DupExpression:
		val := stack.Peek()
		if val == nil {
			return RuntimeErrorf("tried to duplicate top of empty stack: '%s'", e.stackLabel())
		}
		stack.Push(val)
	case ast.SwapExpression:
		if stack.Size() < 2 {
			return RuntimeErrorf("tried to swap top of stack with fewer than 2 items: '%s'", e.stackLabel())
		}
		first, err := stack.Pop()
		if err != nil {
			return err
		}
		second, err := stack.Pop()
		if err != nil {
			return err
		}
		stack.Push(first)
		stack.Push(second)
	default:
		return InternalErrorf("unrecognized stack operation %s", expr.String())
	}

	return nil
}

func (e *EvalState) evalBinaryOperator(expr ast.Expression) (err error) {
	activeStack, err := e.ActiveStack()
	if err != nil {
		return
	}

	vals, err := activeStack.PopN(2)
	if err != nil {
		err = RuntimeErrorf("tried to use a binary operator on stack with fewer than 2 items: '%s'", e.stackLabel())
		return
	}

	if len(vals) != 2 {
		err = InternalError{Message: "stack popping logic failure"}
		return
	}

	x, y := vals[0], vals[1]

	if _, ok := expr.(ast.EqualToExpression); ok {
		activeStack.Push(&BooleanValue{Value: x == y})
		return nil
	}

	res, err := y.DoBinaryOperation(x, expr)
	if err != nil {
		activeStack.Push(y)
		activeStack.Push(x)
		return err
	}
	activeStack.Push(res)
	return
}

func (e *EvalState) evalBitwiseOperator(expr ast.Expression) (err error) {
	activeStack, err := e.ActiveStack()
	if err != nil {
		return
	}

	vals, err := activeStack.PopN(2)
	if err != nil {
		err = RuntimeErrorf("tried to use a binary operator on stack with fewer than 2 items: '%s'", e.stackLabel())
		return
	}

	if len(vals) != 2 {
		err = InternalError{Message: "stack popping logic failure"}
		return
	}

	x, y := vals[0], vals[1]
	res, err := y.DoBitwiseOperation(x, expr)
	if err != nil {
		activeStack.Push(y)
		activeStack.Push(x)
		return err
	}
	activeStack.Push(res)
	return
}

func (e *EvalState) evalAssignment(expr ast.Expression) (err error) {
	switch typedExpr := expr.(type) {
	case ast.AssignmentExpression:
		return e.evalBasicAssignment(typedExpr)
	case ast.NAssignmentExpression:
		return e.evalNAssignment(typedExpr)
	default:
		return InternalError{Message: "invalid eval assignment"}
	}
}

func (e *EvalState) evalBasicAssignment(expr ast.AssignmentExpression) (err error) {
	stack, err := e.ActiveStack()
	if err != nil {
		return err
	}
	val, err := stack.Pop()
	if err != nil {
		return RuntimeErrorf("could not assign a value to stack '%s', not enough values on stack '%s'", expr.Identifier, e.stackLabel())
	}

	err = e.pushOrCreate(expr.Identifier, val)
	return
}

func (e *EvalState) evalNAssignment(expr ast.NAssignmentExpression) (err error) {
	stack, err := e.ActiveStack()
	if err != nil {
		return err
	}
	vals, err := stack.PopN(uint16(expr.N))
	if err != nil {
		return RuntimeErrorf("could not assign %d values to stack '%s', not enough values on stack '%s'", expr.N, expr.Identifier, e.stackLabel())
	}

	for i := len(vals) - 1; i >= 0; i-- {
		if err = e.pushOrCreate(expr.Identifier, vals[i]); err != nil {
			return
		}
	}
	return
}

func (e *EvalState) pushOrCreate(stack string, val Value) error {
	if targetStack, ok := e.NamedStacks[stack]; ok {
		targetStack.Push(val)
	} else {
		if e.HasFunction(stack) {
			return RuntimeErrorf("cannot assign to stack '%s', it is a function name", stack)
		}
		newStack := &Stack{}
		newStack.Push(val)
		e.NamedStacks[stack] = newStack
	}
	return nil
}

func (e *EvalState) evalIdentifier(expr ast.IdentifierExpression) (errors []error) {
	if e.HasFunction(expr.Value) {
		errors = e.executeFunction(expr.Value)
		if len(errors) > 0 {
			return
		}
		return []error{}
	}

	identStack, ok := e.NamedStacks[expr.Value]
	if !ok {
		return []error{RuntimeErrorf("could not recognize identifier '%s'", expr.Value)}
	}

	val, err := identStack.Pop()
	if err != nil {
		return []error{RuntimeErrorf("identifier '%s' is an empty stack", expr.Value)}
	}

	e.MainStack.Push(val)
	return
}
