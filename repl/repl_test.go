package repl

import (
	"testing"
)

func TestCanExit(t *testing.T) {
	p := NewTestPrinterExpectations([]string{
		">>> ",
		"Goodbye!\n",
	})
	Run(NewTestInput([]string{
		"exit",
	}), p)
	if err := p.IsValid(); err != nil {
		t.Error(err)
	}
}

func TestItSkipsEmptyLines(t *testing.T) {
	p := NewTestPrinterExpectations([]string{
		">>> ",
		">>> ",
		"Goodbye!\n",
	})
	Run(NewTestInput([]string{
		"",
		"exit",
	}), p)
	if err := p.IsValid(); err != nil {
		t.Error(err)
	}
}

// TODO: Don't want to write parser blocks by hand, bring back when eval works
// func TestCanPrint(t *testing.T) {
// 	p := NewTestPrinterExpectations([]string{
// 		">>> ",
// 		"(INT 4)\n",
// 		"(FLOAT .20)\n",
// 		"(ADD +)\n",
// 		">>> ",
// 		">>> ",
// 		"(HASH #)\n",
// 		"Illegal token: $\n",
// 		">>> ",
// 		"Goodbye!\n",
// 	})
// 	Run(NewTestInput([]string{
// 		"4 .20 +",
// 		" ",
// 		"#$",
// 		"exit",
// 	}), p)
// 	if err := p.IsValid(); err != nil {
// 		t.Error(err)
// 	}
// }
