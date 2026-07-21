package rd

import (
	"llewellyn-kevin/yak/ast"
	"llewellyn-kevin/yak/lexer"
)

// --------------------------------------------------
// Recursive Descent (Top-Down) Parser Implementation
// --------------------------------------------------

const RECURSIVE_DESCENT_STRATEGY = "recursive-descent"

type RdParser struct {
	l             *lexer.Lexer
	currentToken  lexer.Token
	peekToken     lexer.Token
	blockId       int
	conditionalId int
}

func (RdParser) Strategy() string { return RECURSIVE_DESCENT_STRATEGY }

func NewRdParser(l *lexer.Lexer) *RdParser {
	p := &RdParser{l: l}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *RdParser) Parse() *ast.Program {
	program := &ast.Program{
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

func (p *RdParser) parseExpression() ast.Expression {
	var expr ast.Expression

	switch true {
	case p.expectCurrent(lexer.INT):
		expr = p.parseInt()
	case p.expectCurrent(lexer.FLOAT):
		expr = p.parseFloat()
	case p.expectCurrent(lexer.STRING): // TODO: Not implemented in lexer
		expr = p.parseString()
	case p.expectCurrent(lexer.SYMBOL):
		expr = p.parseSymbol()
	case p.expectCurrent(lexer.TRUE):
		expr = p.parseBool(true)
	case p.expectCurrent(lexer.FALSE):
		expr = p.parseBool(false)
	case p.expectCurrent(lexer.ADD):
		expr = p.parseAdd()
	case p.expectCurrent(lexer.SUB):
		expr = p.parseSubtract()
	case p.expectCurrent(lexer.MULT):
		expr = p.parseMultiply()
	case p.expectCurrent(lexer.DIV):
		expr = p.parseDivide()
	case p.expectCurrent(lexer.MOD):
		expr = p.parseModulo()
	case p.expectCurrent(lexer.EXP):
		expr = p.parsePower()
	case p.expectCurrent(lexer.LT):
		expr = p.parseLessThan()
	case p.expectCurrent(lexer.GT):
		expr = p.parseGreaterThan()
	case p.expectCurrent(lexer.LTEQ):
		expr = p.parseLessThanEqualTo()
	case p.expectCurrent(lexer.GTEQ):
		expr = p.parseGreaterThanEqualTo()
	case p.expectCurrent(lexer.EQ):
		expr = p.parseEqualTo()
	case p.expectCurrent(lexer.SWAP):
		expr = p.parseSwap()
	case p.expectCurrent(lexer.BOR):
		expr = p.parseBor()
	case p.expectCurrent(lexer.BAND):
		expr = p.parseBand()
	case p.expectCurrent(lexer.BXOR):
		expr = p.parseXor()
	case p.expectCurrent(lexer.DUP):
		expr = p.parseDup()
	case p.expectCurrent(lexer.INC):
		expr = p.parseIncrement()
	case p.expectCurrent(lexer.DEC):
		expr = p.parseDecrement()
	case p.expectCurrent(lexer.LSHIFT):
		expr = p.parseLeftShift()
	case p.expectCurrent(lexer.RSHIFT):
		expr = p.parseRightShift()
	case p.expectCurrent(lexer.ASSIGN):
		expr = p.parseAssignment()
	case p.expectCurrent(lexer.IDENT):
		expr = p.parseIdent()
	case p.expectCurrent(lexer.YAKOUT):
		expr = p.parseYakout()
	default:
		return nil
	}

	if p.expectPeek(lexer.NASSIGN) {
		intLit, ok := expr.(ast.IntLiteral)
		if !ok {
			panic("n-assignment requires an integer literal count")
		}
		p.nextToken()
		return p.parseNAssignment(intLit.Value, intLit.Token)
	}

	return expr
}
