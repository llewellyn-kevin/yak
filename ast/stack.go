package ast

import "llewellyn-kevin/yak/lexer"

type SwapExpression struct {
	Token lexer.Token
}

func (SwapExpression) isExpression() {}

func (s SwapExpression) String() string {
	return "<>"
}

type DupExpression struct {
	Token lexer.Token
}

func (DupExpression) isExpression() {}

func (d DupExpression) String() string {
	return "."
}
