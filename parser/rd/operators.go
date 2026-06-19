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

type ModuloExpression struct {
	Token lexer.Token
}

func (ModuloExpression) isExpression() {}

func (m ModuloExpression) String() string {
	return "%"
}

func (p *RdParser) parseModulo() Expression {
	return ModuloExpression{
		Token: p.currentToken,
	}
}

type PowerExpression struct {
	Token lexer.Token
}

func (PowerExpression) isExpression() {}

func (p PowerExpression) String() string {
	return "**"
}

func (p *RdParser) parsePower() Expression {
	return PowerExpression{
		Token: p.currentToken,
	}
}

type LessThanExpression struct {
	Token lexer.Token
}

func (LessThanExpression) isExpression() {}

func (l LessThanExpression) String() string {
	return "<"
}

func (p *RdParser) parseLessThan() Expression {
	return LessThanExpression{
		Token: p.currentToken,
	}
}

type GreaterThanExpression struct {
	Token lexer.Token
}

func (GreaterThanExpression) isExpression() {}

func (g GreaterThanExpression) String() string {
	return ">"
}

func (p *RdParser) parseGreaterThan() Expression {
	return GreaterThanExpression{
		Token: p.currentToken,
	}
}

type LessThanEqualToExpression struct {
	Token lexer.Token
}

func (LessThanEqualToExpression) isExpression() {}

func (l LessThanEqualToExpression) String() string {
	return "<="
}

func (p *RdParser) parseLessThanEqualTo() Expression {
	return LessThanEqualToExpression{
		Token: p.currentToken,
	}
}

type GreaterThanEqualToExpression struct {
	Token lexer.Token
}

func (GreaterThanEqualToExpression) isExpression() {}

func (g GreaterThanEqualToExpression) String() string {
	return ">="
}

func (p *RdParser) parseGreaterThanEqualTo() Expression {
	return GreaterThanEqualToExpression{
		Token: p.currentToken,
	}
}

type EqualToExpression struct {
	Token lexer.Token
}

func (EqualToExpression) isExpression() {}

func (e EqualToExpression) String() string {
	return "="
}

func (p *RdParser) parseEqualTo() Expression {
	return EqualToExpression{
		Token: p.currentToken,
	}
}

type BorExpression struct {
	Token lexer.Token
}

func (BorExpression) isExpression() {}

func (b BorExpression) String() string {
	return "|"
}

func (p *RdParser) parseBor() Expression {
	return BorExpression{
		Token: p.currentToken,
	}
}

type BandExpression struct {
	Token lexer.Token
}

func (BandExpression) isExpression() {}

func (b BandExpression) String() string {
	return "&"
}

func (p *RdParser) parseBand() Expression {
	return BandExpression{
		Token: p.currentToken,
	}
}

type XorExpression struct {
	Token lexer.Token
}

func (XorExpression) isExpression() {}

func (x XorExpression) String() string {
	return "^"
}

func (p *RdParser) parseXor() Expression {
	return XorExpression{
		Token: p.currentToken,
	}
}

type IncrementExpression struct {
	Token lexer.Token
}

func (IncrementExpression) isExpression() {}

func (i IncrementExpression) String() string {
	return "++"
}

func (p *RdParser) parseIncrement() Expression {
	return IncrementExpression{
		Token: p.currentToken,
	}
}

type DecrementExpression struct {
	Token lexer.Token
}

func (DecrementExpression) isExpression() {}

func (d DecrementExpression) String() string {
	return "--"
}

func (p *RdParser) parseDecrement() Expression {
	return DecrementExpression{
		Token: p.currentToken,
	}
}

type LeftShiftExpression struct {
	Token lexer.Token
}

func (LeftShiftExpression) isExpression() {}

func (l LeftShiftExpression) String() string {
	return "<<"
}

func (p *RdParser) parseLeftShift() Expression {
	return LeftShiftExpression{
		Token: p.currentToken,
	}
}

type RightShiftExpression struct {
	Token lexer.Token
}

func (RightShiftExpression) isExpression() {}

func (r RightShiftExpression) String() string {
	return ">>"
}

func (p *RdParser) parseRightShift() Expression {
	return RightShiftExpression{
		Token: p.currentToken,
	}
}
