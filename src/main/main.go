package main

import (
  "fmt"
)

func main() {
  testin := []string{"1", "2", "+", "EOF"}
  for lexer := Newlexer(testin); lexer.Gettoken().Value != "EOF"; lexer.Next() {
    fmt.Println(lexer.Gettoken())
  }
}
