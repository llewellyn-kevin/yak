package parser_test

import (
	"llewellyn-kevin/yak/lexer"
	"llewellyn-kevin/yak/parser"
	"llewellyn-kevin/yak/parser/rd"
	"strings"
	"testing"
)

func TestParserFactory(t *testing.T) {
	l := lexer.NewLexer(strings.NewReader(""))
	actual, err := parser.NewParserFactory().Get(l)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if actual.Strategy() != rd.RECURSIVE_DESCENT_STRATEGY {
		t.Fatalf("Expected strategy %v, got %v", rd.RECURSIVE_DESCENT_STRATEGY, actual.Strategy())
	}
}
