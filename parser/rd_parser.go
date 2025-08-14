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
		Nodes: []Node{},
	}

	for p.currentToken.Type != lexer.EOF {
		node := p.parseStatement()
		if node != nil {
			program.Nodes = append(program.Nodes, node)
		}
		p.nextToken()
	}

	return program
}

func (p *rdParser) parseStatement() Node {
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
	case p.expectCurrent(lexer.ADD):
		return p.parseBinaryOperand(ADD_OPERATOR)
	case p.expectCurrent(lexer.SUB):
		return p.parseBinaryOperand(SUBTRACT_OPERATOR)
	case p.expectCurrent(lexer.MULT):
		return p.parseBinaryOperand(MULTIPLY_OPERATOR)
	case p.expectCurrent(lexer.DIV):
		return p.parseBinaryOperand(DIVIDE_OPERATOR)
	case p.expectCurrent(lexer.MOD):
		return p.parseBinaryOperand(MODULO_OPERATOR)
	case p.expectCurrent(lexer.LT):
		return p.parseBinaryOperand(LT_OPERATOR)
	case p.expectCurrent(lexer.GT):
		return p.parseBinaryOperand(GT_OPERATOR)
	case p.expectCurrent(lexer.LTEQ):
		return p.parseBinaryOperand(LTE_OPERATOR)
	case p.expectCurrent(lexer.GTEQ):
		return p.parseBinaryOperand(GTE_OPERATOR)
	case p.expectCurrent(lexer.EQ):
		return p.parseBinaryOperand(EQUAL_OPERATOR)
	case p.expectCurrent(lexer.SWAP):
		return p.parseBinaryOperand(SWAP_OPERATOR)
	case p.expectCurrent(lexer.BOR):
		return p.parseBinaryOperand(BINARY_OR_OPERATOR)
	case p.expectCurrent(lexer.BAND):
		return p.parseBinaryOperand(BINARY_AND_OPERATOR)
	case p.expectCurrent(lexer.BXOR):
		return p.parseBinaryOperand(BINARY_XOR_OPERATOR)
	case p.expectCurrent(lexer.LSHIFT):
		return p.parseBinaryOperand(LEFT_SHIFT_OPERATOR)
	case p.expectCurrent(lexer.RSHIFT):
		return p.parseBinaryOperand(RIGHT_SHIFT_OPERATOR)
	case p.expectCurrent(lexer.DUP):
		return p.parseUnaryOperand(DUPLICATE_OPERATOR)
	case p.expectCurrent(lexer.INC):
		return p.parseUnaryOperand(INCREMENT_OPERATOR)
	case p.expectCurrent(lexer.DEC):
		return p.parseUnaryOperand(DECREMENT_OPERATOR)
	case p.expectCurrent(lexer.YAKOUT):
		return p.parseUnaryOperand(YAKOUT_OPERATOR)
	case p.expectCurrent(lexer.YAKIN):
		return p.parseUnaryOperand(YAKIN_OPERATOR)
	case p.expectCurrent(lexer.LITERAL_YAKOUT):
		return p.parseUnaryOperand(YAKOUT_LITERAL_OPERATOR)
	case p.expectCurrent(lexer.YAKUP):
		return p.parseUnaryOperand(YAKUP_OPERATOR)
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

func (p *rdParser) parseBinaryOperand(o BinaryOperatorType) Expression {
	return NewBinaryOperator(p.currentToken, o)
}

func (p *rdParser) parseUnaryOperand(o UnaryOperatorType) Expression {
	return NewUnaryOperator(p.currentToken, o)
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
