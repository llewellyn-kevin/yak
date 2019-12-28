package main

import "regexp"

type Token struct {
  Symbol string
  Value string

}

func tokenlookup(value string) string {
  integer := regexp.MustCompile(`^\d*$`)
  operator := regexp.MustCompile(`^\+$|^\-$|^\*$|^/$`)
  eof := regexp.MustCompile(`^EOF$`)

  switch {
  case integer.FindString(value) != "":
    return "INTEGER"
  case operator.FindString(value) != "":
    return "OPERATOR"
  case eof.FindString(value) != "":
    return "EOF"
  default:
    return "UNKNOWN"
  }
}

func Newtoken(s string) *Token {
  t := new(Token)
  t.Value = s
  t.Symbol = tokenlookup(s)
  return t
}

func (t Token) String() string {
  return "<" + t.Symbol + ", " + t.Value + ">"
}
