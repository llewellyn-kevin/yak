package lexer

import "fmt"

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

func (t Token) String() string {
	return fmt.Sprintf("(%s %s)", t.Type, t.Literal)
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers
	IDENT = "IDENT"

	// Literals
	INT    = "INT"
	FLOAT  = "FLOAT"
	STRING = "STRING" // Not Implemented Yet
	SYMBOL = "SYMBOL"

	// Operators
	ADD  = "ADD"
	SUB  = "SUB"
	MULT = "MULT"
	DIV  = "DIV"
	MOD  = "MOD"
	LT   = "LT"
	GT   = "GT"
	LTEQ = "LTEQ"
	GTEQ = "GTEQ"
	EQ   = "EQ"
	SWAP = "SWAP"
	BOR  = "BOR"
	BAND = "BAND"
	BXOR = "BXOR"

	DUP    = "DUP"
	INC    = "INC"
	DEC    = "DEC"
	LSHIFT = "LSHIFT"
	RSHIFT = "RSHIFT"

	// Delimiters
	LBRACE  = "LBRACE"
	RBRACE  = "RBRACE"
	LPAREN  = "LPAREN"
	RPARENT = "RPAREN"
	HASH    = "HASH"
	COLON   = "COLON"
	SQUOTE  = "SQUOTE"
	DQUOTE  = "DQUOTE"

	// Control
	IF   = "IF"
	NOT  = "NOT"
	ELSE = "ELSE"
	FOR  = "FOR"

	// Assignments
	ASSIGN  = "ASSIGN"
	NASSIGN = "NASSIGN" // Not Implemented Yet

	// Keywords
	TRUE    = "TRUE"
	FALSE   = "FALSE"
	KBOOL   = "KBOOL"
	KINT    = "INT"
	KFLOAT  = "KFLOAT"
	KSTRING = "KSTRING"
	KSYMBOL = "KSYMBOL"

	SET     = "SET"
	SETG    = "SETG"
	SETOPTS = "SETOPTS"
	GET     = "GET"

	YAKOUT         = "YAKOUT"
	YAKIN          = "YAKIN"
	LITERAL_YAKOUT = "LITERAL_YAKOUT"
	YAKUP          = "YAKUP"
)
