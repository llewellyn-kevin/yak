package main

import (
  "fmt"
  "os"
  "bytes"
  "io/ioutil"
)

const (
  BUFFER_END = "<EOF>"
)

func main() {
  // Grab filename to compile from command line args
  fname := os.Args[1]

  // Grab the file and read into a string
  file, err := ioutil.ReadFile(fname)
  if err != nil {
    panic(err)
  }

  // Generate a string array to pass to the lexer, add an EOF
  input := getfields(file)
  input = append(input, BUFFER_END)

  // Create a syntax tree and lexer, run through parser
  tree := Newtree(*Newtoken("<MAIN>"))
  lexer := Newlexer(input)
  Parse(lexer, tree.Root, BUFFER_END)

  fmt.Println(tree)
}

/**
 * Takes a byte array and converts it to a string array ignoring whitespace
 */
func getfields(bytearr []byte) []string {
  const SPACE byte = 32
  const TAB byte = 9
  const BREAK byte = 10

  buffempty := true
  stringbuff := bytes.NewBufferString("")
  stringarr := []string{}

  for _, value := range bytearr {
    // TODO: Fix the switch statement so code is repeated 3 times
    switch value {
    case SPACE:
      if !buffempty {
        stringarr = append(stringarr, stringbuff.String())
        stringbuff = bytes.NewBufferString("")
        buffempty = true
      }
    case TAB:
      if !buffempty {
        stringarr = append(stringarr, stringbuff.String())
        stringbuff = bytes.NewBufferString("")
        buffempty = true
      }
    case BREAK:
      if !buffempty {
        stringarr = append(stringarr, stringbuff.String())
        stringbuff = bytes.NewBufferString("")
        buffempty = true
      }
    default:
      stringbuff.WriteByte(value)
      buffempty = false
    }
  }

  return stringarr
}
