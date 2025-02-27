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
	offset int
}

func NewLexer(input ReadSeeker) *Lexer {
	return &Lexer{
		Input:  input,
		offset: 0,
	}
}

func (l *Lexer) NextToken() (output Token) {
	var buffer []byte

	for output.Type == "" {
		nextChar := make([]byte, 1)
		n, err := l.Input.Read(nextChar)
		l.offset += n

		if err == io.EOF {
			output = Token{EOF, ""}
			return
		}

		buffer = append(buffer, nextChar...)
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

			switch peeked {
			case byte('-'):
				output = Token{DEC, "--"}
				advancer()
			case byte('>'):
				output = Token{ASSIGN, "->"}
				advancer()
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

			// Single / Identifiers
		case ".": // . or .IDENT
			output = Token{DUP, string(buffer)}
		}
	}

	return
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
