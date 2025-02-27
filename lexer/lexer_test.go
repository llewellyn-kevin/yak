package lexer

import (
	"strings"
	"testing"
)

func TestNextToken(t *testing.T) {
	reader := strings.NewReader("%")
	lexer := NewLexer(reader)

	expected := []Token{{MOD, "%"}}

	for _, token := range expected {
		if actual := lexer.NextToken(); actual != token {
			t.Fatalf("Unexpected token. Want %v, Got %v", token, actual)
		}
	}
}

func TestSingleCharOps(t *testing.T) {
	reader := strings.NewReader("+-*/%=&|^")
	lexer := NewLexer(reader)

	expected := []Token{
		{ADD, "+"},
		{SUB, "-"},
		{MULT, "*"},
		{DIV, "/"},
		{MOD, "%"},
		{EQ, "="},
		{BAND, "&"},
		{BOR, "|"},
		{BXOR, "^"},
	}

	for _, token := range expected {
		if actual := lexer.NextToken(); actual != token {
			t.Fatalf("Unexpected token. Want %v, Got %v", token, actual)
		}
	}
}

func TestDoubleCharOps(t *testing.T) {
	reader := strings.NewReader("><++--+->-<<<>>><=>=")
	lexer := NewLexer(reader)

	expected := []Token{
		{GT, ">"},
		{LT, "<"},
		{INC, "++"},
		{DEC, "--"},
		{ADD, "+"},
		{ASSIGN, "->"},
		{SUB, "-"},
		{LSHIFT, "<<"},
		{SWAP, "<>"},
		{RSHIFT, ">>"},
		{LTEQ, "<="},
		{GTEQ, ">="},
	}

	for _, token := range expected {
		if actual := lexer.NextToken(); actual != token {
			t.Fatalf("Unexpected token. Want %v, Got %v", token, actual)
		}
	}
}
