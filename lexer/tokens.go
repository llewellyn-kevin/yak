package lexer

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers
	IDENT  = "IDENT"
	INT    = "INT"
	FLOAT  = "FLOAT"
	STRING = "STRING"
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

	// ?
	HYPHEN = "HYPHEN"

	// Control
	IF   = "IF"
	NOT  = "NOT"
	ELSE = "ELSE"
	FOR  = "FOR"

	// Keywords
	ASSIGN  = "ASSIGN"
	NASSIGN = "NASSIGN"

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
	YAKINN         = "YAKIN"
	LITERAL_YAKOUT = "LITERAL_YAKOUT"
	YAKUP          = "YAKUP"
)
