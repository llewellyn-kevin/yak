package ast

import "llewellyn-kevin/yak/lexer"

type YakoutExpression struct {
	Token lexer.Token
}

func (YakoutExpression) isExpression() {}

func (e YakoutExpression) String() string {
	return "yakout"
}
