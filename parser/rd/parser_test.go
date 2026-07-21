package rd_test

import (
	"fmt"
	"llewellyn-kevin/yak/lexer"
	"llewellyn-kevin/yak/parser/rd"
	"strings"
	"testing"
)

func TestParserWorks(t *testing.T) {
	program := ``
	expected := `block (
    id: 0
    statements: [
    ]
    expressions: [
    ]
)`
	testParserOutput(t, program, expected)
}

func TestParsingMainBlock(t *testing.T) {
	program := `1 2 +`
	expected := `block (
    id: 0
    statements: [
    ]
    expressions: [
        1
        2
        +
    ]
)`
	testParserOutput(t, program, expected)
}

func TestParsingNestedBlocks(t *testing.T) {
	program := `0 { 2 { 3 4 } { + } + } +`
	expected := `block (
    id: 0
    statements: [
        block (
            id: 1
            statements: [
                block (
                    id: 2
                    statements: [
                    ]
                    expressions: [
                        3
                        4
                    ]
                )
                block (
                    id: 3
                    statements: [
                    ]
                    expressions: [
                        +
                    ]
                )
            ]
            expressions: [
                2
                execute-block 2
                execute-block 3
                +
            ]
        )
    ]
    expressions: [
        0
        execute-block 1
        +
    ]
)`
	testParserOutput(t, program, expected)
}

func TestParsingExpressions(t *testing.T) {
	program := `1 4.2 %sym true false + - * / % ** < > <= >= = <> | & ^ . ++ -- << >> ident`
	expected := `block (
    id: 0
    statements: [
    ]
    expressions: [
        1
        4.2
        %sym
        true
        false
        +
        -
        *
        /
        %
        **
        <
        >
        <=
        >=
        =
        <>
        |
        &
        ^
        .
        ++
        --
        <<
        >>
        ident
    ]
)`
	testParserOutput(t, program, expected)
}

func TestParsingNamedScopeBlocks(t *testing.T) {
	program := `{ + }
    { myScope: + }
    { myForScope for: + }
    { myForScopeWithoutColon for + }
    { notAScope + }
    { for + }
    { each - }`
	expected := `block (
    id: 0
    statements: [
        block (
            id: 1
            statements: [
            ]
            expressions: [
                +
            ]
        )
        block (
            id: 2
            scope: myScope
            statements: [
            ]
            expressions: [
                +
            ]
        )
        block (
            id: 3
            scope: myForScope
            loop: true
            statements: [
            ]
            expressions: [
                +
            ]
        )
        block (
            id: 4
            scope: myForScopeWithoutColon
            loop: true
            statements: [
            ]
            expressions: [
                +
            ]
        )
        block (
            id: 5
            statements: [
            ]
            expressions: [
                notAScope
                +
            ]
        )
        block (
            id: 6
            loop: true
            statements: [
            ]
            expressions: [
                +
            ]
        )
        block (
            id: 7
            loop: true
            statements: [
            ]
            expressions: [
                -
            ]
        )
    ]
    expressions: [
        execute-block 1
        execute-block 2
        execute-block 3
        execute-block 4
        execute-block 5
        execute-block 6
        execute-block 7
    ]
)`
	testParserOutput(t, program, expected)
}

func TestParsingBasicIfStatement(t *testing.T) {
	program := `if { + }`
	expected := `block (
    id: 0
    statements: [
        conditional (
            id: 0
            when-true:
                block (
                    id: 1
                    statements: [
                    ]
                    expressions: [
                        +
                    ]
                )
            when-false:
                block (
                    id: 2
                    statements: [
                    ]
                    expressions: [
                    ]
                )
        )
    ]
    expressions: [
        execute-conditional 0
    ]
)`
	testParserOutput(t, program, expected)
}

func TestParsingIfElseStatement(t *testing.T) {
	program := `3 4 if { 1 + } else { 2 - } *`
	expected := `block (
    id: 0
    statements: [
        conditional (
            id: 0
            when-true:
                block (
                    id: 1
                    statements: [
                    ]
                    expressions: [
                        1
                        +
                    ]
                )
            when-false:
                block (
                    id: 2
                    statements: [
                    ]
                    expressions: [
                        2
                        -
                    ]
                )
        )
    ]
    expressions: [
        3
        4
        execute-conditional 0
        *
    ]
)`
	testParserOutput(t, program, expected)
}

func TestParsingNotStatement(t *testing.T) {

	program := `not { 1 }`
	expected := `block (
    id: 0
    statements: [
        conditional (
            id: 0
            inverted: true
            when-true:
                block (
                    id: 1
                    statements: [
                    ]
                    expressions: [
                        1
                    ]
                )
            when-false:
                block (
                    id: 2
                    statements: [
                    ]
                    expressions: [
                    ]
                )
        )
    ]
    expressions: [
        execute-conditional 0
    ]
)`
	testParserOutput(t, program, expected)
}

func TestParsingAssignment(t *testing.T) {
	program := `42->foo`
	expected := `block (
    id: 0
    statements: [
    ]
    expressions: [
        42
        assign foo
    ]
)`
	testParserOutput(t, program, expected)
}

func TestParsingNAssignment(t *testing.T) {
	program := `42 24 2 =>foo`
	expected := `block (
    id: 0
    statements: [
    ]
    expressions: [
        42
        24
        assign-2 foo
    ]
)`
	testParserOutput(t, program, expected)
}

func TestStio(t *testing.T) {
	program := `yakout`
	expected := `block (
    id: 0
    statements: [
    ]
    expressions: [
        yakout
    ]
)`
	testParserOutput(t, program, expected)
}

func getParser(reader *strings.Reader) *rd.RdParser {
	return rd.NewRdParser(lexer.NewLexer(reader))
}

func testParserOutput(t *testing.T, program string, expected string) {
	ast := getParser(strings.NewReader(program)).Parse()
	stringCompare(t, ast, expected)
}

func stringCompare(t *testing.T, actual fmt.Stringer, expected string) {
	actualStr := actual.String()
	normalized := strings.ReplaceAll(actualStr, "\t", "    ")
	if normalized != expected {
		t.Fatalf(`Unexpected string formatting. Have:
%s
Want:
%s`, actualStr, expected)
	}
}
