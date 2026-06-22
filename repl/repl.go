package repl

import (
	"bufio"
	"fmt"
	"io"
	"llewellyn-kevin/yak/lexer"
	"llewellyn-kevin/yak/parser"
	"llewellyn-kevin/yak/parser/rd"
	"strings"
)

const EXIT_KEY = "exit"

type Printer interface {
	Println(s string)
	Print(s string)
}

// Start the REPL shell.
func Run(in io.Reader, p Printer) {
	for {
		line, runAgain := Prompt(in, p)
		if !runAgain {
			break
		}
		if line == "" {
			continue
		}
		l := lexer.NewLexer(strings.NewReader(line))
		p, err := parser.NewParserFactory().Use(rd.RECURSIVE_DESCENT_STRATEGY).Get(l)
		if err != nil {
			fmt.Println(err)
			break
		}
		program := p.Parse()
		fmt.Println(program)
	}
	p.Println("Goodbye!")
}

// Prompts the user for input and returns true if they did not request to exit.
func Prompt(in io.Reader, p Printer) (string, bool) {
	reader := bufio.NewReader(in)
	p.Print(">>> ")
	text, err := reader.ReadString('\n')
	text = strings.TrimSuffix(text, "\n")
	if err != nil {
		p.Println("Bad input")
	}
	return text, text != EXIT_KEY
}
