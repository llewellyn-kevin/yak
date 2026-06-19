package rd

import (
	"llewellyn-kevin/yak/lexer"
	"llewellyn-kevin/yak/parser"
)

// --------------------------------------------------
// Recursive Descent (Top-Down) Parser Implementation
// --------------------------------------------------

type RdParser struct {
	l            *lexer.Lexer
	currentToken lexer.Token
	peekToken    lexer.Token
	blockId      int
}

func (RdParser) Strategy() parser.ParsingStrategy { return parser.RECURSIVE_DESCENT_STRATEGY }

func NewRdParser(l *lexer.Lexer) *RdParser {
	p := &RdParser{l: l}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *RdParser) Parse() *Program {
	program := &Program{
		MainBlock: p.parseBlock(0),
	}

	return program
}

func (p RdParser) expectCurrent(t lexer.TokenType) bool {
	return p.currentToken.Type == t
}

func (p RdParser) expectPeek(t lexer.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *RdParser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *RdParser) parseExpression() Expression {
	switch true {
	case p.expectCurrent(lexer.INT):
		return p.parseInt()
	case p.expectCurrent(lexer.FLOAT):
		return p.parseFloat()
	case p.expectCurrent(lexer.STRING): // TODO: Not implemented in lexer
		return p.parseString()
	case p.expectCurrent(lexer.SYMBOL):
		return p.parseSymbol()
	case p.expectCurrent(lexer.ADD):
		return p.parseAdd()
	case p.expectCurrent(lexer.SUB):
		return p.parseSubtract()
	case p.expectCurrent(lexer.MULT):
		return p.parseMultiply()
	case p.expectCurrent(lexer.DIV):
		return p.parseDivide()
	case p.expectCurrent(lexer.IDENT):
		return p.parseIdent()
	default:
		return nil
	}
}
