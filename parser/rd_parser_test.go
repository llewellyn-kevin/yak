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
  identifier: 
  {
    token: identifier
    value: foo
    type: any
  }
}
{
  token: assignment
  identifier:
  {
    token: identifier
    value: bar
    type: any
  }
}
{
  token: assignment
  identifier:
  {
    token: identifier
    value: three
    type: any
  }
}
{
  token: assignment
  identifier:
  {
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

func TestRdOperands(t *testing.T) {
	p, err := parser.NewParserFactory().
		Use(parser.RECURSIVE_DESCENT_STRATEGY).
		Get(lexer.NewLexer(strings.NewReader(`+ - * / % < > <= >= = <> | & ^ << >>
. ++ -- yakout yakin yakout! yakup`)))
	program := ensureValidProgram(t, 23, p, err)
	actual := program.String()

	expected := `{
  token: binary-operator
  type: add
}
{
  token: binary-operator
  type: subtract
}
{
  token: binary-operator
  type: multiply
}
{
  token: binary-operator
  type: divide
}
{
  token: binary-operator
  type: mod
}
{
  token: binary-operator
  type: less-than
}
{
  token: binary-operator
  type: greater-than
}
{
  token: binary-operator
  type: less-than-or-equal
}
{
  token: binary-operator
  type: greater-than-or-equal
}
{
  token: binary-operator
  type: equal
}
{
  token: binary-operator
  type: swap
}
{
  token: binary-operator
  type: binary-or
}
{
  token: binary-operator
  type: binary-and
}
{
  token: binary-operator
  type: binary-xor
}
{
  token: binary-operator
  type: left-shift
}
{
  token: binary-operator
  type: right-shift
}
{
  token: unary-operator
  type: duplicate
}
{
  token: unary-operator
  type: increment
}
{
  token: unary-operator
  type: decrement
}
{
  token: unary-operator
  type: yakout
}
{
  token: unary-operator
  type: yakin
}
{
  token: unary-operator
  type: yakout-literal
}
{
  token: unary-operator
  type: yakup
}`

	checkMultiString(t, actual, expected)
}

func TestRdConditional(t *testing.T) {
	p, err := parser.NewParserFactory().
		Use(parser.RECURSIVE_DESCENT_STRATEGY).
		Get(lexer.NewLexer(strings.NewReader(`
1 . -> foo =
if {
    3
    if {
      2 +
    }
} else {
    1
    +
}
not {
    3
}
`)))
	program := ensureValidProgram(t, 6, p, err)
	actual := program.String()

	expected := `{
  token: literal
  type: integer-literal
  value: 1
}
{
  token: unary-operator
  type: duplicate
}
{
  token: assignment
  identifier: 
  {
    token: identifier
    value: foo
    type: any
  }
}
{
  token: binary-operator
  type: equal
}
{
  token: conditional-expression
  type: if
  consequence: [
    {
      token: literal
      type: integer-literal
      value: 3
    }
    {
      token: conditional-expression
      type: if
      consequence: [
        {
          token: literal
          type: integer-literal
          value: 2
        }
        {
          token: binary-operator
          type: add
        }
      ]
      alternative: [
      ]
    }
  ]
  alternative: [
    {
      token: literal
      type: integer-literal
      value: 1
    }
    {
      token: binary-operator
      type: add
    }
  ]
}
{
  token: conditional-expression
  type: not
  consequence: [
    {
      token: literal
      type: integer-literal
      value: 3
    }
  ]
  alternative: [
  ]
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
  identifier:
  {
    token: identifier
    value: nums
    type:
    {
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

	if c := len(program.Nodes); c != targetLength {
		t.Fatalf("program.Statements does not have the correct number of statements. Have %d, want %d.", c, targetLength)
	}

	return program
}
