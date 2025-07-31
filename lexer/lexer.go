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

func (l *Lexer) nextChar() (byte, error) {
	nextChar := make([]byte, 1)
	n, err := l.Input.Read(nextChar)
	l.setOffset(l.offset + int64(n))
	return nextChar[0], err
}

func (l *Lexer) nextWithoutWhitespace() (byte, error) {
	nextChar := byte(' ')
	var err error
	for (nextChar == ' ' || nextChar == '\t' || nextChar == '\n' || nextChar == '\r') && err == nil {
		nextChar, err = l.nextChar()
	}
	return nextChar, err
}

func (l *Lexer) peekChar() (bool, byte, func() (int64, error)) {
	advancer := func() (int64, error) {
		l.offset += 1
		return l.Input.Seek(int64(l.offset), 0)
	}

	char := make([]byte, 1)
	l.Input.Read(char)
	_, err := l.Input.Seek(int64(l.offset), 0)
	if err == io.EOF {
		return true, byte(0), advancer
	}
	return false, char[0], advancer
}

var Keywords = map[string]TokenType{
	"if":      IF,
	"not":     NOT,
	"else":    ELSE,
	"for":     FOR,
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

func (l *Lexer) NextToken() (output Token) {
	var buffer []byte

	for output.Type == "" {
		nextChar, err := l.nextWithoutWhitespace()

		if err == io.EOF {
			output = Token{EOF, ""}
			return
		}

		buffer = append(buffer, nextChar)
		switch string(buffer) {
		// Single Chars
		case "*":
			output = Token{MULT, string(buffer)}
		case "/":
			output = Token{DIV, string(buffer)}
		case "%":
			output = Token{MOD, string(buffer)}
		case "=":
			output = Token{EQ, string(buffer)}
		case "^":
			output = Token{BXOR, string(buffer)}
		case "&":
			output = Token{BAND, string(buffer)}
		case "|":
			output = Token{BOR, string(buffer)}

		// Delimiters
		case "#":
			output = Token{HASH, string(buffer)}
		case "{":
			output = Token{LBRACE, string(buffer)}
		case "}":
			output = Token{RBRACE, string(buffer)}
		case "(":
			output = Token{LPAREN, string(buffer)}
		case ")":
			output = Token{RPARENT, string(buffer)}
		case "'":
			output = Token{SQUOTE, string(buffer)}
		case "\"":
			output = Token{DQUOTE, string(buffer)}

		case ".":
			_, peeked, _ := l.peekChar()
			if isDigit(peeked) {
				output = l.readNumber([]byte{'.'})
			} else {
				output = Token{DUP, string(buffer)}
			}

		case ":":
			_, peeked, _ := l.peekChar()
			if isLetter(peeked) {
				output = l.readSymbol()
			} else {
				output = Token{COLON, string(buffer)}
			}

			// Single / Double
		case "+": // + or ++
			if eof, peeked, advancer := l.peekChar(); eof || peeked != byte('+') {
				output = Token{ADD, "+"}
			} else {
				output = Token{INC, "++"}
				advancer()
			}
		case "-": // - or -- or ->
			eof, peeked, advancer := l.peekChar()
			if eof {
				output = Token{SUB, "-"}
			}

			switch true {
			case peeked == byte('-'):
				output = Token{DEC, "--"}
				advancer()
			case peeked == byte('>'):
				output = Token{ASSIGN, "->"}
				advancer()
			case isDigit(peeked) || peeked == byte('.'):
				output = l.readNumber([]byte{'-'})
			default:
				output = Token{SUB, "-"}
			}
		case "<": // < or <= or <> or <<
			eof, peeked, advancer := l.peekChar()
			if eof {
				output = Token{LT, "<"}
			}

			switch peeked {
			case byte('='):
				output = Token{LTEQ, "<="}
				advancer()
			case byte('<'):
				output = Token{LSHIFT, "<<"}
				advancer()
			case byte('>'):
				output = Token{SWAP, "<>"}
				advancer()
			default:
				output = Token{LT, "<"}
			}
		case ">": // > or >= or >>
			eof, peeked, advancer := l.peekChar()
			if eof {
				output = Token{GT, ">"}
			}

			switch peeked {
			case byte('='):
				output = Token{GTEQ, ">="}
				advancer()
			case byte('>'):
				output = Token{RSHIFT, ">>"}
				advancer()
			default:
				output = Token{GT, ">"}
			}

		// Keywords
		default:
			if isLetter(buffer[0]) {
				output = l.readIdentifier()
			} else if isDigit(buffer[0]) {
				output = l.readNumber(buffer)
			} else {
				output = Token{ILLEGAL, string(buffer)}
			}
		}
	}

	return
}

func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func lookupIdentifier(ident string) TokenType {
	if tok, ok := Keywords[ident]; ok {
		return tok
	}
	return IDENT
}

func (l *Lexer) setOffset(offset int64) {
	if offset < 0 {
		offset = 0
	}
	l.offset = offset
	l.Input.Seek(l.offset, io.SeekStart)
}

func (l *Lexer) goBack(amount int64) {
	l.setOffset(l.offset - amount)
}

func (l *Lexer) readIdentifier() Token {
	var buffer []byte
	l.goBack(1)
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
		} else if err == io.EOF || (!isLetter(nextChar[0]) && !isDigit(nextChar[0]) && nextChar[0] != '_') {
			l.goBack(int64(n))
			break
		}

		buffer = append(buffer, nextChar...)
	}

	s := string(buffer)
	return Token{lookupIdentifier(string(s)), string(s)}
}

func (l *Lexer) readNumber(buffer []byte) Token {
	isFloat := false
	if buffer[0] == '.' {
		isFloat = true
	}

	for {
		n, err := l.nextChar()

		if err == io.EOF {
			break
		}

		if n == byte('.') {
			if isFloat {
				l.goBack(1)
				break
			}
			isFloat = true
		} else if !isDigit(n) {
			l.goBack(1)
			break
		}

		buffer = append(buffer, n)
	}
	if isFloat {
		return Token{FLOAT, string(buffer)}
	}
	return Token{INT, string(buffer)}
}

func (l *Lexer) readSymbol() Token {
	buffer := []byte{':'}
	for {
		n, err := l.nextChar()

		if err == io.EOF {
			break
		}

		if !isLetter(n) && !isDigit(n) && n != '-' {
			l.goBack(1)
			break
		}

		buffer = append(buffer, n)
	}

	return Token{SYMBOL, string(buffer)}
}
