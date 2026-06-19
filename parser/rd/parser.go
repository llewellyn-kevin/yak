package rd

import (
	"llewellyn-kevin/yak/lexer"
	"llewellyn-kevin/yak/parser"
)

// --------------------------------------------------
// Recursive Descent (Top-Down) Parser Implementation
// --------------------------------------------------

type RdParser struct {
	l             *lexer.Lexer
	currentToken  lexer.Token
	peekToken     lexer.Token
	blockId       int
	conditionalId int
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
	case p.expectCurrent(lexer.TRUE):
		return p.parseBool(true)
	case p.expectCurrent(lexer.FALSE):
		return p.parseBool(false)
	case p.expectCurrent(lexer.ADD):
		return p.parseAdd()
	case p.expectCurrent(lexer.SUB):
		return p.parseSubtract()
	case p.expectCurrent(lexer.MULT):
		return p.parseMultiply()
	case p.expectCurrent(lexer.DIV):
		return p.parseDivide()
	case p.expectCurrent(lexer.MOD):
		return p.parseModulo()
	case p.expectCurrent(lexer.EXP):
		return p.parsePower()
	case p.expectCurrent(lexer.LT):
		return p.parseLessThan()
	case p.expectCurrent(lexer.GT):
		return p.parseGreaterThan()
	case p.expectCurrent(lexer.LTEQ):
		return p.parseLessThanEqualTo()
	case p.expectCurrent(lexer.GTEQ):
		return p.parseGreaterThanEqualTo()
	case p.expectCurrent(lexer.EQ):
		return p.parseEqualTo()
	case p.expectCurrent(lexer.SWAP):
		return p.parseSwap()
	case p.expectCurrent(lexer.BOR):
		return p.parseBor()
	case p.expectCurrent(lexer.BAND):
		return p.parseBand()
	case p.expectCurrent(lexer.BXOR):
		return p.parseXor()
	case p.expectCurrent(lexer.DUP):
		return p.parseDup()
	case p.expectCurrent(lexer.INC):
		return p.parseIncrement()
	case p.expectCurrent(lexer.DEC):
		return p.parseDecrement()
	case p.expectCurrent(lexer.LSHIFT):
		return p.parseLeftShift()
	case p.expectCurrent(lexer.RSHIFT):
		return p.parseRightShift()
	case p.expectCurrent(lexer.ASSIGN):
		return p.parseAssignment()
	case p.expectCurrent(lexer.IDENT):
		return p.parseIdent()
	case p.expectCurrent(lexer.YAKOUT):
		return p.parseYakout()
	default:
		return nil
	}
}
