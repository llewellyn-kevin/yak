package main

import (
	"fmt"
	"llewellyn-kevin/yak/repl"
	"os"
)

type StdoutPrinter struct{}

func (StdoutPrinter) Println(s string) {
	fmt.Println(s)
}
func (StdoutPrinter) Print(s string) {
	fmt.Print(s)
}

func main() {
	repl.Run(os.Stdin, &StdoutPrinter{})
}
