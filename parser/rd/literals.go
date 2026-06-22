package rd

import (
	"llewellyn-kevin/yak/ast"
	"llewellyn-kevin/yak/lexer"
	"strconv"
	"strings"
)

func (p *RdParser) parseInt() ast.Expression {
	i, err := strconv.Atoi(p.currentToken.Literal)
	if err != nil {
		// TODO: Handle, but this should never happen since the lexer should only return valid int literals
		panic(err)
	}

	return ast.IntLiteral{
		Token: p.currentToken,
		Value: i,
	}
}

func (p *RdParser) parseFloat() ast.Expression {
	f, err := strconv.ParseFloat(p.currentToken.Literal, 64)
	if err != nil {
		panic(err)
	}

	return ast.FloatLiteral{
		Token: p.currentToken,
		Value: f,
	}
}

func (p *RdParser) parseString() ast.Expression {
	return ast.StringLiteral{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}
}

func (p *RdParser) parseSymbol() ast.Expression {
	return ast.SymbolLiteral{
		Token: p.currentToken,
		Value: strings.Replace(p.currentToken.Literal, "%", "", 1),
	}
}

func (p *RdParser) parseBool(value bool) ast.Expression {
	return ast.BoolLiteral{
		Token: p.currentToken,
		Value: value,
	}
}

func (p *RdParser) parseIdent() ast.Expression {
	return ast.IdentifierExpression{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}
}

func (p *RdParser) parseAssignment() ast.Expression {
	assignmentToken := p.currentToken
	p.nextToken()

	if !p.expectCurrent(lexer.IDENT) {
		// TODO: Add to more robust error handler
		panic("Expected identifier after assignment operator")
	}

	return ast.AssignmentExpression{
		Tokens: []lexer.Token{
			assignmentToken,
			p.currentToken,
		},
		Identifier: p.currentToken.Literal,
	}
}
