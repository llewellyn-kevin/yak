package lexer

import (
	"io"
)

type ReadSeeker interface {
	io.Reader
	io.Seeker
}

type Lexer struct {
	Input  ReadSeeker
	offset int64
}

func NewLexer(input ReadSeeker) *Lexer {
	return &Lexer{
		Input:  input,
		offset: 0,
	}
}

// A map of byte arrays representing keywords to their correct token.
var Keywords = map[string]TokenType{
	"if":      IF,
	"not":     NOT,
	"else":    ELSE,
	"for":     FOR,
	"each":    FOR,
	"bool":    KBOOL,
	"int":     KINT,
	"float":   KFLOAT,
	"string":  KSTRING,
	"symbol":  KSYMBOL,
	"true":    TRUE,
	"false":   FALSE,
	"set":     SET,
	"setg":    SETG,
	"setopts": SETOPTS,
	"get":     GET,
	"yakout":  YAKOUT,
	"yakin":   YAKIN,
	"yakup":   YAKUP,
	"yakout!": LITERAL_YAKOUT,
}

// Maps byte arrays representing tokens to their token outputers.
var tokenMap = tMap{
	"*":  terminalToken(MULT),
	"/":  terminalToken(DIV),
	"=":  terminalToken(EQ),
	"^":  terminalToken(BXOR),
	"&":  terminalToken(BAND),
	"|":  terminalToken(BOR),
	"#":  terminalToken(HASH),
	":":  terminalToken(COLON),
	"{":  terminalToken(LBRACE),
	"}":  terminalToken(RBRACE),
	"(":  terminalToken(LPAREN),
	")":  terminalToken(RPARENT),
	"'":  terminalToken(SQUOTE),
	"\"": terminalToken(DQUOTE),
	"+": recursiveToken(terminalToken(ADD), tMap{
		"+": terminalToken(INC),
	}),
	"<": recursiveToken(terminalToken(LT), tMap{
		"=": terminalToken(LTEQ),
		">": terminalToken(SWAP),
		"<": terminalToken(LSHIFT),
	}),
	">": recursiveToken(terminalToken(GT), tMap{
		"=": terminalToken(GTEQ),
		">": terminalToken(RSHIFT),
	}),
	".": recursiveToken(
		terminalToken(DUP),
		tMap{},
		logicalHandler{
			func(peek []byte) bool {
				return isDigit(peek)
			},
			func(ch []byte, l *Lexer) Token {
				return l.readNumber(ch)
			},
		},
	),
	"%": recursiveToken(
		terminalToken(MOD),
		tMap{},
		logicalHandler{
			func(peek []byte) bool {
				return isLetter(peek)
			},
			func(ch []byte, l *Lexer) Token {
				return l.readSymbol(ch)
			},
		},
	),
	"-": recursiveToken(
		terminalToken(SUB),
		tMap{
			"-": terminalToken(DEC),
			">": terminalToken(ASSIGN),
		},
		logicalHandler{
			func(peek []byte) bool {
				return isDigit(peek) || (len(peek) == len([]byte{'.'}) && peek[0] == '.')
			},
			func(ch []byte, l *Lexer) Token {
				return l.readNumber(ch)
			},
		},
	),
}

// A function that takes a char array and returns a token.
type tokenOutputter func(ch []byte, l *Lexer) Token

// Maps characters to their tokens via token outputters.
type tMap map[string]tokenOutputter

// Returns a closure that returns the token for any symbol that is only a
// single character or is on the last character, so no lookahead is needed.
func terminalToken(t TokenType) tokenOutputter {
	return func(ch []byte, _ *Lexer) Token {
		return Token{t, string(ch)}
	}
}

// Struct used to hold logical checks where matching a token is more complex
// than matching a single character.
type logicalHandler struct {
	check   func(peek []byte) bool
	scanner tokenOutputter
}

// Returns a closure to handle reading multi-character tokens.
// fallback is the token to return if no match is found in the map.
// nextMap is the map of functions that return tokens with the symbols as
// their keys.
// lHandlers are an optional list of logicalHandlers. If the check function
// for the handler returns true, the scanner function will be executed to get
// the token. Otherwise this will be skipped. This can be used for checks where
// multiple characters are valid for the given token.
func recursiveToken(fallback tokenOutputter, nextMap tMap, lHandlers ...logicalHandler) tokenOutputter {
	return func(ch []byte, l *Lexer) Token {
		n, err := l.nextChar()
		if err == io.EOF {
			return fallback(ch, l)
		}

		out, exists := nextMap[string(n)]
		if exists {
			ch = append(ch, n...)
			return out(ch, l)
		}

		for _, h := range lHandlers {
			if h.check(n) {
				ch = append(ch, n...)
				return h.scanner(ch, l)
			}
		}

		l.goBack(1)
		return fallback(ch, l)
	}
}

// Get the next semantically relevant token.
func (l *Lexer) NextToken() (output Token) {
	for output.Type == "" {
		nextChar, err := l.nextWithoutWhitespace()

		if err == io.EOF {
			output = Token{EOF, ""}
			return
		}

		t, exists := tokenMap[string(nextChar)]
		if exists {
			output = t(nextChar, l)
			continue
		}

		if isLetter(nextChar) {
			output = l.readIdentifier(nextChar)
		} else if isDigit(nextChar) {
			output = l.readNumber(nextChar)
		} else {
			output = Token{ILLEGAL, string(nextChar)}
		}
	}

	return
}

// Lexer helper that returns the next rune as a byte array
func (l *Lexer) nextChar() ([]byte, error) {
	nextChar := make([]byte, 1)
	n, err := l.Input.Read(nextChar)
	l.setOffset(l.offset + int64(n))
	return nextChar, err
}

// Lexer helper that gets the next char but skips any whitespace
func (l *Lexer) nextWithoutWhitespace() ([]byte, error) {
	nextChar := []byte{' '}
	var err error
	for len(nextChar) == 1 && (nextChar[0] == ' ' || nextChar[0] == '\t' || nextChar[0] == '\n' || nextChar[0] == '\r') && err == nil {
		nextChar, err = l.nextChar()
	}
	return nextChar, err
}

// Lexer helper that jumps the file reader to a specific offset
func (l *Lexer) setOffset(offset int64) {
	if offset < 0 {
		offset = 0
	}
	l.offset = offset
	l.Input.Seek(l.offset, io.SeekStart)
}

// Lexer helper that goes back the specified number of bytes in the reader
func (l *Lexer) goBack(amount int64) {
	l.setOffset(l.offset - amount)
}

// Returns whether the rune contained in the byte array is a letter or underscore
func isLetter(ch []byte) bool {
	if len(ch) > len([]byte{'a'}) {
		return false
	}
	b := ch[0]
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || b == '_'
}

// Returns whether or not the rune contained in the byte array is a number
func isDigit(ch []byte) bool {
	if len(ch) > len([]byte{'0'}) {
		return false
	}
	b := ch[0]
	return b >= '0' && b <= '9'
}

// See if the given string is a language keyword. And return that keyword's
// token if so. Otherwise, return the identifier token with the string as
// content.
func lookupIdentifier(ident string) TokenType {
	if tok, ok := Keywords[ident]; ok {
		return tok
	}
	return IDENT
}

// Keep reading in symbols in the lexer until you hit a non valid id char.
// Return the resulting keyword or identifier token.
func (l *Lexer) readIdentifier(buffer []byte) Token {
	for {
		nextChar := make([]byte, 1)
		n, err := l.Input.Read(nextChar)
		l.offset += int64(n)

		if nextChar[0] == '!' { // Handle special case for yakout!
			if len(buffer) > 0 && string(buffer) == "yakout" {
				return Token{LITERAL_YAKOUT, "yakout!"}
			}
			l.goBack(int64(n))
			break
		} else if err == io.EOF || (!isLetter(nextChar) && !isDigit(nextChar) && nextChar[0] != '_') {
			l.goBack(int64(n))
			break
		}

		buffer = append(buffer, nextChar...)
	}

	s := string(buffer)
	return Token{lookupIdentifier(string(s)), string(s)}
}

// Keep reading the file until an invalid input for a number. Return the correct
// number literal token (INT or FLOAT)
func (l *Lexer) readNumber(buffer []byte) Token {
	isFloat := false
	if buffer[0] == '.' || (buffer[0] == '-' && buffer[1] == '.') {
		isFloat = true
	}

	for {
		n, err := l.nextChar()

		if err == io.EOF {
			break
		}

		if n[0] == '.' {
			if isFloat {
				l.goBack(1)
				break
			}
			isFloat = true
		} else if !isDigit(n) {
			l.goBack(1)
			break
		}

		buffer = append(buffer, n...)
	}
	if isFloat {
		return Token{FLOAT, string(buffer)}
	}
	return Token{INT, string(buffer)}
}

// Keep reading the file until an invalid symbol input. Then return the symbol
// token with the read content.
func (l *Lexer) readSymbol(buffer []byte) Token {
	for {
		n, err := l.nextChar()

		if err == io.EOF {
			break
		}

		if !isLetter(n) && !isDigit(n) && n[0] != '-' {
			l.goBack(1)
			break
		}

		buffer = append(buffer, n...)
	}

	return Token{SYMBOL, string(buffer)}
}
