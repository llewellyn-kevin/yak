package rd

import "llewellyn-kevin/yak/lexer"

type AddExpression struct {
	Token lexer.Token
}

func (AddExpression) isExpression() {}

func (a AddExpression) String() string {
	return "+"
}

func (p *RdParser) parseAdd() Expression {
	return AddExpression{
		Token: p.currentToken,
	}
}

type SubtractExpression struct {
	Token lexer.Token
}

func (SubtractExpression) isExpression() {}

func (s SubtractExpression) String() string {
	return "-"
}

func (p *RdParser) parseSubtract() Expression {
	return SubtractExpression{
		Token: p.currentToken,
	}
}

type MultiplyExpression struct {
	Token lexer.Token
}

func (MultiplyExpression) isExpression() {}

func (m MultiplyExpression) String() string {
	return "*"
}

func (p *RdParser) parseMultiply() Expression {
	return MultiplyExpression{
		Token: p.currentToken,
	}
}

type DivideExpression struct {
	Token lexer.Token
}

func (DivideExpression) isExpression() {}

func (d DivideExpression) String() string {
	return "/"
}

func (p *RdParser) parseDivide() Expression {
	return DivideExpression{
		Token: p.currentToken,
	}
}
