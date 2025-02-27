package lexer

import "testing"

func TestNextToken(t *testing.T) {
	lexer := Lexer{}

	expected := []Token{
		{ILLEGAL, "foobar"},
	}

	for _, token := range expected {
		if actual := lexer.NextToken(); actual != token {
			t.Fatalf("Unexpected token. Want %v, Got %v", token, actual)
		}
	}
}
