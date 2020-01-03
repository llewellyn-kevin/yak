package main

import "regexp"

type Token struct {
  Symbol string
  Value string

}

func tokenlookup(value string) string {
  main := regexp.MustCompile(`^<MAIN>$`)
  integer := regexp.MustCompile(`^\d+$`)
  float := regexp.MustCompile(`^\d+\.\d+$`)
  operator := regexp.MustCompile(`^\+$|^\-$|^\*$|^/$`)
  identifier := regexp.MustCompile(`^[_\$a-zA-Z][_\a-zA-Z0-9]*$`)
  eof := regexp.MustCompile(`^<EOF>$`)

  switch {
  case main.FindString(value) != "":
    return "MAIN"
  case integer.FindString(value) != "":
    return "INTEGER"
  case float.FindString(value) != "":
    return "FLOAT"
  case operator.FindString(value) != "":
    return "OPERATOR"
  case identifier.FindString(value) != "":
    return "IDENTIFIER"
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
