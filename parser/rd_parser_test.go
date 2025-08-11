package parser_test

import (
	"llewellyn-kevin/yak/lexer"
	"llewellyn-kevin/yak/parser"
	"strings"
	"testing"
)

func TestRdAssignment(t *testing.T) {
	p, err := parser.NewParserFactory().
		Use(parser.RECURSIVE_DESCENT_STRATEGY).
		Get(lexer.NewLexer(strings.NewReader(`-> foo
-> bar ->three->four`)))
	program := ensureValidProgram(t, 4, p, err)
	actual := program.String()

	expected := `{
  token: assignment
  identifier: {
    token: identifier
    value: foo
    type: any
  }
}
{
  token: assignment
  identifier: {
    token: identifier
    value: bar
    type: any
  }
}
{
  token: assignment
  identifier: {
    token: identifier
    value: three
    type: any
  }
}
{
  token: assignment
  identifier: {
    token: identifier
    value: four
    type: any
  }
}`

	checkMultiString(t, actual, expected)
}

func TestRdLiteral(t *testing.T) {
	p, err := parser.NewParserFactory().
		Use(parser.RECURSIVE_DESCENT_STRATEGY).
		Get(lexer.NewLexer(strings.NewReader(`1 2.3 true false %sym`)))
	program := ensureValidProgram(t, 5, p, err)
	actual := program.String()

	expected := `{
  token: literal
  type: integer-literal
  value: 1
}
{
  token: literal
  type: float-literal
  value: 2.3
}
{
  token: literal
  type: bool-literal
  value: true
}
{
  token: literal
  type: bool-literal
  value: false
}
{
  token: literal
  type: symbol-literal
  value: %sym
}`

	checkMultiString(t, actual, expected)
}

func TestRdTypedAssignment(t *testing.T) {
	p, err := parser.NewParserFactory().
		Use(parser.RECURSIVE_DESCENT_STRATEGY).
		Get(lexer.NewLexer(strings.NewReader(`1 -> nums:int|string|sym`)))
	program := ensureValidProgram(t, 2, p, err)
	actual := program.String()

	expected := `{
  token: literal
  type: integer-literal
  value: 1
}
{
  token: assignment
  identifier: {
    token: identifier
    value: nums
    type: {
      token: type-declaration
      types: [
        {
          token: keyword-type
          value: int
        }
        {
          token: keyword-type
          value: string
        }
        {
          token: symbolic-type
          value: sym
        }
      ]
    }
  }
}`

	checkMultiString(t, actual, expected)
}

func ensureValidProgram(t *testing.T, targetLength int, p parser.Parser, err error) *parser.Program {
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	program := p.Parse()
	if program == nil {
		t.Fatalf("Parse() returned nil")
	}

	if c := len(program.Statements); c != targetLength {
		t.Fatalf("program.Statements does not have the correct number of statements. Have %d, want %d.", c, targetLength)
	}

	return program
}
