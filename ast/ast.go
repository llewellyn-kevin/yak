package ast

type Statement interface {
	isStatement()
	String() string
}

type Expression interface {
	isExpression()
	String() string
}
