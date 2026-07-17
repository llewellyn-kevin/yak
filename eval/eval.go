package eval

import (
	"errors"
	"fmt"
	"llewellyn-kevin/yak/ast"
)

func Eval(p *ast.Program) (*EvalState, []error) {
	state := NewEvalState()
	errors := state.executeBlock(p.MainBlock)
	return state, errors
}

func PartialEval(p *ast.Program, state *EvalState) []error {
	if state == nil {
		state = NewEvalState()
	}
	errors := state.executeBlock(p.MainBlock)
	return errors
}

func (e *EvalState) executeBlock(block *ast.Block) (errors []error) {
	for _, expr := range block.Expressions {
		var err error
		if blockExp, ok := expr.(ast.ExecuteBlockExpression); ok {
			block, err = block.GetBlock(blockExp.BlockId)
			if err == nil {
				errors = append(errors, e.executeBlock(block)...)
			}
		} else {
			err = e.evalExpression(expr)
		}

		if err != nil {
			errors = append(errors, err)
		}
	}
	return
}

func (e *EvalState) evalExpression(expr ast.Expression) (err error) {
	switch {
	case ast.IsLiteral(expr):
		err = e.evalLiteral(expr)
	case ast.IsUnaryOperator(expr):
		err = e.evalUnaryOperator(expr)
	case ast.IsBinaryOperator(expr):
		err = e.evalBinaryOperator(expr)
	case ast.IsBitwiseOperator(expr):
		err = e.evalBitwiseOperator(expr)
	default:
		err = fmt.Errorf("Runtime Error: unknown expression %s", expr.String())
	}
	return
}

func (e *EvalState) evalLiteral(expr ast.Expression) (err error) {
	switch v := expr.(type) {
	case ast.IntLiteral:
		e.MainStack.Push(IntValue{Value: v.Value})
	case ast.FloatLiteral:
		e.MainStack.Push(FloatValue{Value: v.Value})
	case ast.StringLiteral:
		e.MainStack.Push(StringValue{Value: v.Value})
	case ast.SymbolLiteral:
		e.MainStack.Push(SymbolValue{Value: v.Value})
	case ast.BoolLiteral:
		e.MainStack.Push(BooleanValue{Value: v.Value})
	default:
		err = fmt.Errorf("Runtime Error: unknown literal type %s", expr.String())
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
		err = fmt.Errorf("Runtime Error: tried to use a unary operator on empty stack: '%s'", e.stackLabel())
		return
	}

	result, err := literal.DoUnaryOperation(expr)
	activeStack.Push(result)
	return
}

func (e *EvalState) evalBinaryOperator(expr ast.Expression) (err error) {
	activeStack, err := e.ActiveStack()
	if err != nil {
		return
	}

	vals, err := activeStack.PopN(2)
	if err != nil {
		err = fmt.Errorf("Runtime Error: tried to use a binary operator on stack with fewer than 2 items: '%s'", e.stackLabel())
		return
	}

	if len(vals) != 2 {
		err = errors.New("Internal Error: stack popping logic failure")
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
		err = fmt.Errorf("Runtime Error: tried to use a binary operator on stack with fewer than 2 items: '%s'", e.stackLabel())
		return
	}

	if len(vals) != 2 {
		err = errors.New("Internal Error: stack popping logic failure")
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
