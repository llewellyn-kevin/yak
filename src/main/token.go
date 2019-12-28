package main

type Token struct {
  Symbol string
  Value string
}

func Newtoken(s string) *Token {
  t := new(Token)
  t.Symbol = s
  t.Value = s
  return t
}

func (t Token) String() string {
  return "<" + t.Symbol + ", " + t.Value + ">"
}
