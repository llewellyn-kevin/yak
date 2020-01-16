package main

type Lexer struct {
  Fields []string
  Numfields int
  Scanpos int
}

func Newlexer(fields []string) *Lexer {
  l := new(Lexer)
  l.Fields = fields
  l.Scanpos = 0
  l.Numfields = len(l.Fields)
  return l
}

func (l Lexer) getfield() string {
  return l.Fields[l.Scanpos]
}

func (l *Lexer) Next() {
  if l.Scanpos == l.Numfields - 1 {
    l.Scanpos = 0
  } else {
    l.Scanpos = l.Scanpos + 1
  }
}

func (l Lexer) Gettoken() *Token {
  return Newtoken(l.getfield())
}
