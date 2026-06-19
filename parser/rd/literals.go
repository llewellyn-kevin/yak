package rd

import (
	"llewellyn-kevin/yak/lexer"
	"strconv"
	"strings"
)

type IntLiteral struct {
	Token lexer.Token
	Value int
}

func (IntLiteral) isExpression() {}

func (i IntLiteral) String() string {
	return strconv.Itoa(i.Value)
}

func (p *RdParser) parseInt() Expression {
	i, err := strconv.Atoi(p.currentToken.Literal)
	if err != nil {
		// TODO: Handle, but this should never happen since the lexer should only return valid int literals
		panic(err)
	}

	return IntLiteral{
		Token: p.currentToken,
		Value: i,
	}
}

type FloatLiteral struct {
	Token lexer.Token
	Value float64
}

func (FloatLiteral) isExpression() {}

func (f FloatLiteral) String() string {
	return strconv.FormatFloat(f.Value, 'f', -1, 64)
}

func (p *RdParser) parseFloat() Expression {
	f, err := strconv.ParseFloat(p.currentToken.Literal, 64)
	if err != nil {
		panic(err)
	}

	return FloatLiteral{
		Token: p.currentToken,
		Value: f,
	}
}

type StringLiteral struct {
	Token lexer.Token
	Value string
}

func (StringLiteral) isExpression() {}

func (s StringLiteral) String() string {
	return strconv.Quote(s.Value)
}

func (p *RdParser) parseString() Expression {
	return StringLiteral{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}
}

type SymbolLiteral struct {
	Token lexer.Token
	Value string
}

func (SymbolLiteral) isExpression() {}

func (s SymbolLiteral) String() string {
	return "%" + s.Value
}

func (p *RdParser) parseSymbol() Expression {
	return SymbolLiteral{
		Token: p.currentToken,
		Value: strings.Replace(p.currentToken.Literal, "%", "", 1),
	}
}

type BoolLiteral struct {
	Token lexer.Token
	Value bool
}

func (BoolLiteral) isExpression() {}

func (b BoolLiteral) String() string {
	return strconv.FormatBool(b.Value)
}

func (p *RdParser) parseBool(value bool) Expression {
	return BoolLiteral{
		Token: p.currentToken,
		Value: value,
	}
}

type IdentifierExpression struct {
	Token lexer.Token
	Value string
}

func (IdentifierExpression) isExpression() {}

func (i IdentifierExpression) String() string {
	return i.Value
}

func (p *RdParser) parseIdent() Expression {
	return IdentifierExpression{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}
}
