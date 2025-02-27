package lexer

type Lexer struct {
}

func (l *Lexer) NextToken() Token {
	return Token{ILLEGAL, "foobar"}
}
