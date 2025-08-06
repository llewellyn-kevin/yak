package repl

import (
	"bufio"
	"io"
	"llewellyn-kevin/yak/lexer"
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
		for t := l.NextToken(); t.Type != lexer.EOF; t = l.NextToken() {
			if t.Type == lexer.ILLEGAL {
				p.Println("Illegal token: " + t.Literal)
				continue
			}
			p.Println(t.String())
		}
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
