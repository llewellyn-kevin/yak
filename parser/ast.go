package parser

import (
	"fmt"
	"llewellyn-kevin/yak/lexer"
	"strings"
)

type Node interface {
	RecStringer
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type RecStringer interface {
	String(nestLevel int) string
}

type Program struct {
	Statements []Statement
}

// Represents an assignment operation with its identifier
type Assignment struct {
	Token      lexer.Token
	Identifier *Identifier
}

type ValidType string

const (
	BOOL_TYPE   = "bool"
	INT_TYPE    = "int"
	FLOAT_TYPE  = "float"
	STRING_TYPE = "string"
	SYMBOL_TYPE = "symbol"
)

type TypeDef interface {
	RecStringer
	isTypeExpression()
}

// Represents a set of valid types
type TypeDeclaration struct {
	Token lexer.Token
	Types []TypeDef
}

// Represents a valid primitive type in the language
type KeywordType struct {
	Token lexer.Token
	Value ValidType
}

// Represents a valid symbol name that can be used as a value in the field
type SymbolicType struct {
	Token lexer.Token
	Value string
}

type LiteralType string

const (
	BOOL_LITERAL   = "bool-literal"
	INT_LITERAL    = "integer-literal"
	FLOAT_LITERAL  = "float-literal"
	STRING_LITERAL = "string-literal"
	SYMBOL_LITERAL = "symbol-literal"
)

type Literal struct {
	Token lexer.Token
	Type  LiteralType
	Value string
}

type Identifier struct {
	Token   lexer.Token
	Value   string
	TypeSet TypeDeclaration
}

// ---------------------------------------------------------
// Program Implementation
// ---------------------------------------------------------
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p Program) String() (o string) {
	for i, s := range p.Statements {
		o += s.String(0)
		if i < (len(p.Statements) - 1) {
			o += "\n"
		}
	}
	return
}

// ---------------------------------------------------------
// helpers
// ---------------------------------------------------------
func nestedLines(lines []string, nestLevel int) (o string) {
	padding := strings.Repeat("  ", nestLevel)
	for i, l := range lines {
		if i == 0 {
			o += fmt.Sprintf("%s", l)
		} else {
			o += fmt.Sprintf("%s%s", padding, l)
		}
		if i != (len(lines) - 1) {
			o += "\n"
		}
	}
	return
}

// ---------------------------------------------------------
// Assignment Implementation
// ---------------------------------------------------------
func (Assignment) statementNode()         {}
func (a Assignment) TokenLiteral() string { return a.Token.Literal }
func (a Assignment) String(nestLevel int) string {
	return nestedLines([]string{
		"{",
		"  token: assignment",
		fmt.Sprintf("  identifier: %s", a.Identifier.String(nestLevel+1)),
		"}",
	}, nestLevel)
}
func NewAssignment(a lexer.Token, i lexer.Token, typeDef *TypeDeclaration) *Assignment {
	return &Assignment{
		Token:      a,
		Identifier: NewIdentifier(i, typeDef),
	}
}

// ---------------------------------------------------------
// Types Implementation
// ---------------------------------------------------------
func (TypeDeclaration) expressionNode()        {}
func (t TypeDeclaration) TokenLiteral() string { return t.Token.Literal }
func (t TypeDeclaration) String(nestLevel int) string {
	var typeList []string
	for _, t := range t.Types {
		typeList = append(typeList, t.String(nestLevel+2))
	}
	return nestedLines(append(append([]string{
		"{",
		"  token: type-declaration",
		"  types: [",
	}, typeList...), []string{
		"  ]",
		"}",
	}...), nestLevel)
}
func (KeywordType) expressionNode()        {}
func (KeywordType) isTypeExpression()      {}
func (t KeywordType) TokenLiteral() string { return t.Token.Literal }
func (t KeywordType) String(nestLevel int) string {
	padding := strings.Repeat("  ", nestLevel-2)
	return nestedLines([]string{
		fmt.Sprintf("%s{", padding),
		"  token: keyword-type",
		fmt.Sprintf("  value: %s", t.Value),
		"}",
	}, nestLevel)
}
func (SymbolicType) expressionNode()        {}
func (SymbolicType) isTypeExpression()      {}
func (t SymbolicType) TokenLiteral() string { return t.Token.Literal }
func (t SymbolicType) String(nestLevel int) string {
	padding := strings.Repeat("  ", nestLevel-2)
	return nestedLines([]string{
		fmt.Sprintf("%s{", padding),
		"  token: symbolic-type",
		fmt.Sprintf("  value: %s", t.Value),
		"}",
	}, nestLevel)
}

// ---------------------------------------------------------
// Identifier Implementation
// ---------------------------------------------------------
func (Identifier) statementNode()         {}
func (i Identifier) TokenLiteral() string { return i.Token.Literal }
func (i Identifier) String(nestLevel int) string {
	var types string
	if len(i.TypeSet.Types) == 0 {
		types = "any"
	} else {
		types = i.TypeSet.String(nestLevel + 1)
	}
	return nestedLines([]string{
		"{",
		"  token: identifier",
		fmt.Sprintf("  value: %s", i.Value),
		fmt.Sprintf("  type: %s", types),
		"}",
	}, nestLevel)
}
func NewIdentifier(i lexer.Token, typeDef *TypeDeclaration) *Identifier {
	return &Identifier{
		Token:   i,
		Value:   i.Literal,
		TypeSet: *typeDef,
	}
}

// ---------------------------------------------------------
// Literal Implementation
// ---------------------------------------------------------
func (Literal) statementNode()         {}
func (i Literal) TokenLiteral() string { return i.Token.Literal }
func (i Literal) String(nestLevel int) string {
	return nestedLines([]string{
		"{",
		"  token: literal",
		fmt.Sprintf("  type: %s", i.Type),
		fmt.Sprintf("  value: %s", i.Value),
		"}",
	}, nestLevel)
}
func NewLiteral(i lexer.Token, t LiteralType) *Literal {
	return &Literal{
		Token: i,
		Type:  t,
		Value: i.Literal,
	}
}
