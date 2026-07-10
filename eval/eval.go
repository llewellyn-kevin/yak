package eval

import (
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
		err := e.evalExpression(expr)
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
	return err
}
