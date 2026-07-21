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
	EXP  = "EXP"
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
	SQUOTE  = "SQUOTE" // Not yet in parser
	DQUOTE  = "DQUOTE" // Not yet in parser

	// Control
	IF   = "IF"
	NOT  = "NOT"
	ELSE = "ELSE"
	FOR  = "FOR"

	// Assignments
	ASSIGN  = "ASSIGN"
	NASSIGN = "NASSIGN"

	// Keywords
	TRUE    = "TRUE"
	FALSE   = "FALSE"
	KBOOL   = "KBOOL"
	KINT    = "KINT"
	KFLOAT  = "KFLOAT"
	KSTRING = "KSTRING"
	KSYMBOL = "KSYMBOL"

	SET     = "SET"     // Not yet in parser
	SETG    = "SETG"    // Not yet in parser
	SETOPTS = "SETOPTS" // Not yet in parser
	GET     = "GET"     // Not yet in parser

	YAKOUT         = "YAKOUT"
	YAKIN          = "YAKIN"          // Not yet in parser
	LITERAL_YAKOUT = "LITERAL_YAKOUT" // Not yet in parser
	YAKUP          = "YAKUP"          // Not yet in parser
)
