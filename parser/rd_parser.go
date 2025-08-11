package parser

import (
	"errors"
	"llewellyn-kevin/yak/lexer"
)

// --------------------------------------------------
// Recursive Descent (Top-Down) Parser Implementation
// --------------------------------------------------

type rdParser struct {
	l *lexer.Lexer

	currentToken lexer.Token
	peekToken    lexer.Token
}

func (p *rdParser) Strategy() ParsingStrategy { return RECURSIVE_DESCENT_STRATEGY }

func (p *rdParser) Parse() *Program {
	program := &Program{
		Statements: []Statement{},
	}

	for p.currentToken.Type != lexer.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *rdParser) parseStatement() Statement {
	switch true {
	case p.expectCurrent(lexer.ASSIGN):
		return p.parseAssignStatement()
	case p.expectCurrent(lexer.TRUE) || p.expectCurrent(lexer.FALSE):
		return p.parseLiteral(BOOL_LITERAL)
	case p.expectCurrent(lexer.INT):
		return p.parseLiteral(INT_LITERAL)
	case p.expectCurrent(lexer.FLOAT):
		return p.parseLiteral(FLOAT_LITERAL)
	case p.expectCurrent(lexer.STRING):
		return p.parseLiteral(STRING_LITERAL)
	case p.expectCurrent(lexer.SYMBOL):
		return p.parseLiteral(SYMBOL_LITERAL)
	default:
		return nil
	}
}

func (p *rdParser) parseTypeList(typeDeclaration *TypeDeclaration) error {
	var t TypeDef
	if p.expectCurrent(lexer.IDENT) {
		t = SymbolicType{p.currentToken, p.currentToken.Literal}
	} else {
		var ty ValidType
		switch p.currentToken.Type {
		case lexer.KBOOL:
			ty = BOOL_TYPE
		case lexer.KINT:
			ty = INT_TYPE
		case lexer.KFLOAT:
			ty = FLOAT_TYPE
		case lexer.KSTRING:
			ty = STRING_TYPE
		case lexer.KSYMBOL:
			ty = SYMBOL_TYPE
		default:
			return errors.New("Invalid token, expected valid type")
		}
		t = KeywordType{p.currentToken, ty}
	}
	typeDeclaration.Types = append(typeDeclaration.Types, t)
	p.nextToken()
	if p.expectCurrent(lexer.BOR) {
		p.nextToken()
		return p.parseTypeList(typeDeclaration)
	}
	return nil
}

func (p *rdParser) parseAssignStatement() Statement {
	if !p.expectPeek(lexer.IDENT) {
		return nil // should error out?
	}
	assign := p.currentToken
	p.nextToken()
	ident := p.currentToken
	typeDef := &TypeDeclaration{}

	if p.expectPeek(lexer.COLON) {
		p.nextToken()
		p.nextToken()
		p.parseTypeList(typeDef)
	}

	return NewAssignment(assign, ident, typeDef)
}

func (p *rdParser) parseLiteral(ty LiteralType) Statement {
	return NewLiteral(p.currentToken, ty)
}

func (p rdParser) expectCurrent(t lexer.TokenType) bool {
	return p.currentToken.Type == t
}

func (p rdParser) expectPeek(t lexer.TokenType) bool {
	return p.peekToken.Type == t
}

func newRdParser(l *lexer.Lexer) *rdParser {
	p := &rdParser{l: l}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *rdParser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.l.NextToken()
}
