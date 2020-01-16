package main

import "regexp"

type Token struct {
  Symbol string
  Value string

}

func tokenlookup(value string) string {
  main := regexp.MustCompile(`^<MAIN>$`)

  checkif := regexp.MustCompile(`^\?$`)
  checknot := regexp.MustCompile(`^!$`)

  blockopen := regexp.MustCompile(`^\{$`)
  blockclose := regexp.MustCompile(`^\}$`)

  integer := regexp.MustCompile(`^\d+$`)
  float := regexp.MustCompile(`^\d+\.\d+$`)
  operator := regexp.MustCompile(`^\+$|^\-$|^\*$|^/$|^%$|^==$|^\.$`)
  identifier := regexp.MustCompile(`^[_\$a-zA-Z][_\a-zA-Z0-9]*$`)
  funcheader := regexp.MustCompile(`^\d+#[_\$a-zA-Z][_\a-zA-Z0-9]*$`)

  eof := regexp.MustCompile(`^<EOF>$`)

  switch {
  case main.FindString(value) != "":
    return "MAIN"
  case checkif.FindString(value) != "":
    return "IF"
  case checknot.FindString(value) != "":
    return "NOT"
  case blockopen.FindString(value) != "":
    return "BLOCK_OPEN"
  case blockclose.FindString(value) != "":
    return "BLOCK_CLOSE"
  case funcheader.FindString(value) != "":
    return "FUNC_HEADER"
  case integer.FindString(value) != "":
    return "INTEGER"
  case float.FindString(value) != "":
    return "FLOAT"
  case operator.FindString(value) != "":
    return "OPERATOR"
  case identifier.FindString(value) != "":
    return "IDENTIFIER"
  case eof.FindString(value) != "":
    return "<EOF>"
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
