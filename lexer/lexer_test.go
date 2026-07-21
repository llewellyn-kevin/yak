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
	reader := strings.NewReader("+-/%=&|^><")
	lexer := NewLexer(reader)

	expected := []Token{
		{ADD, "+"},
		{SUB, "-"},
		{DIV, "/"},
		{MOD, "%"},
		{EQ, "="},
		{BAND, "&"},
		{BOR, "|"},
		{BXOR, "^"},
		{GT, ">"},
		{LT, "<"},
	}

	for _, token := range expected {
		if actual := lexer.NextToken(); actual != token {
			t.Fatalf("Unexpected token. Want %v, Got %v", token, actual)
		}
	}
}

func TestDoubleCharOps(t *testing.T) {
	reader := strings.NewReader("*><++--+->-<<<>>><= >=..**=>")
	lexer := NewLexer(reader)

	expected := []Token{
		{MULT, "*"},
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
		{DUP, "."},
		{DUP, "."},
		{EXP, "**"},
		{NASSIGN, "=>"},
		{EOF, ""},
	}

	for _, token := range expected {
		if actual := lexer.NextToken(); actual != token {
			t.Fatalf("Unexpected token. Want %v, Got %v", token, actual)
		}
	}
}

func TestDelimiters(t *testing.T) {
	reader := strings.NewReader("#{}():|'\"")
	lexer := NewLexer(reader)

	expected := []Token{
		{HASH, "#"},
		{LBRACE, "{"},
		{RBRACE, "}"},
		{LPAREN, "("},
		{RPARENT, ")"},
		{COLON, ":"},
		{BOR, "|"},
		{SQUOTE, "'"},
		{DQUOTE, "\""},
		{EOF, ""},
	}

	for _, token := range expected {
		if actual := lexer.NextToken(); actual != token {
			t.Fatalf("Unexpected token. Want %v, Got %v", token, actual)
		}
	}
}

func TestControlFlow(t *testing.T) {
	reader := strings.NewReader("if not else for")
	lexer := NewLexer(reader)

	expected := []Token{
		{IF, "if"},
		{NOT, "not"},
		{ELSE, "else"},
		{FOR, "for"},
		{EOF, ""},
	}

	for _, token := range expected {
		if actual := lexer.NextToken(); actual != token {
			t.Fatalf("Unexpected token. Want %v, Got %v", token, actual)
		}
	}
}

func TestKeywords(t *testing.T) {
	reader := strings.NewReader("bool int float string symbol true false set setg setopts get yakout yakin yakup yakout!")
	lexer := NewLexer(reader)

	expected := []Token{
		{KBOOL, "bool"},
		{KINT, "int"},
		{KFLOAT, "float"},
		{KSTRING, "string"},
		{KSYMBOL, "symbol"},
		{TRUE, "true"},
		{FALSE, "false"},
		{SET, "set"},
		{SETG, "setg"},
		{SETOPTS, "setopts"},
		{GET, "get"},
		{YAKOUT, "yakout"},
		{YAKIN, "yakin"},
		{YAKUP, "yakup"},
		{LITERAL_YAKOUT, "yakout!"},
		{EOF, ""},
	}

	for _, token := range expected {
		if actual := lexer.NextToken(); actual != token {
			t.Fatalf("Unexpected token. Want %v, Got %v", token, actual)
		}
	}
}

func TestIdentifiers(t *testing.T) {
	reader := strings.NewReader("foo bar_baz qux123 _private")
	lexer := NewLexer(reader)

	expected := []Token{
		{IDENT, "foo"},
		{IDENT, "bar_baz"},
		{IDENT, "qux123"},
		{IDENT, "_private"},
		{EOF, ""},
	}

	for _, token := range expected {
		if actual := lexer.NextToken(); actual != token {
			t.Fatalf("Unexpected token. Want %v, Got %v", token, actual)
		}
	}
}

func TestNumbers(t *testing.T) {
	reader := strings.NewReader("1 1.2 -358 -13.21 .34 0.55 -.89 1.4. 4:233. 377")
	lexer := NewLexer(reader)

	expected := []Token{
		{INT, "1"},
		{FLOAT, "1.2"},
		{INT, "-358"},
		{FLOAT, "-13.21"},
		{FLOAT, ".34"},
		{FLOAT, "0.55"},
		{FLOAT, "-.89"},
		{FLOAT, "1.4"},
		{DUP, "."},
		{INT, "4"},
		{COLON, ":"},
		{FLOAT, "233."},
		{INT, "377"},
		{EOF, ""},
	}

	for _, token := range expected {
		if actual := lexer.NextToken(); actual != token {
			t.Fatalf("Unexpected token. Want %v, Got %v", token, actual)
		}
	}
}

func TestSymbols(t *testing.T) {
	reader := strings.NewReader("%foo %foo-bar %-f %0b %m0m %mt-4 %fizz-buzz-")
	lexer := NewLexer(reader)

	expected := []Token{
		{SYMBOL, "%foo"},
		{SYMBOL, "%foo-bar"},
		{MOD, "%"},
		{SUB, "-"},
		{IDENT, "f"},
		{MOD, "%"},
		{INT, "0"},
		{IDENT, "b"},
		{SYMBOL, "%m0m"},
		{SYMBOL, "%mt-4"},
		{SYMBOL, "%fizz-buzz-"},
		{EOF, ""},
	}

	for _, token := range expected {
		if actual := lexer.NextToken(); actual != token {
			t.Fatalf("Unexpected token. Want %v, Got %v", token, actual)
		}
	}
}

func TestExceptions(t *testing.T) {
	reader := strings.NewReader("!@ foo!")
	lexer := NewLexer(reader)

	expected := []Token{
		{ILLEGAL, "!"},
		{ILLEGAL, "@"},
		{IDENT, "foo"},
		{ILLEGAL, "!"},
		{EOF, ""},
	}

	for _, token := range expected {
		if actual := lexer.NextToken(); actual != token {
			t.Fatalf("Unexpected token. Want %v, Got %v", token, actual)
		}
	}
}

func TestWhitespace(t *testing.T) {
	reader := strings.NewReader(" \t\n\r foo  bar\tbaz ")
	lexer := NewLexer(reader)

	expected := []Token{
		{IDENT, "foo"},
		{IDENT, "bar"},
		{IDENT, "baz"},
		{EOF, ""},
	}

	for _, token := range expected {
		if actual := lexer.NextToken(); actual != token {
			t.Fatalf("Unexpected token. Want %v, Got %v", token, actual)
		}
	}
}
