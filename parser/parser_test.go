package parser_test

import (
	"llewellyn-kevin/yak/lexer"
	"llewellyn-kevin/yak/parser"
	"strings"
	"testing"
)

func TestParserFactory(t *testing.T) {
	l := lexer.NewLexer(strings.NewReader(""))
	actual, err := parser.NewParserFactory().Get(l)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if actual.Strategy() != parser.RECURSIVE_DESCENT_STRATEGY {
		t.Fatalf("Expected strategy %v, got %v", parser.RECURSIVE_DESCENT_STRATEGY, actual.Strategy())
	}
}

func TestRecursiveStrings(t *testing.T) {
	aTok := lexer.Token{Type: lexer.ASSIGN, Literal: "->"}
	iTok := lexer.Token{Type: lexer.IDENT, Literal: "foo"}
	statement := parser.NewAssignment(aTok, iTok, &parser.TypeDeclaration{})
	actual := statement.String(0)

	expected := `{
  token: assignment
  identifier: {
    token: identifier
    value: foo
    type: any
  }
}`

	checkMultiString(t, actual, expected)
}

func checkMultiString(t *testing.T, actual string, expected string) {
	if actual != expected {
		t.Fatalf(`Unexpected string formatting. Have:
%s

Want:

%s`, actual, expected)
	}
}
