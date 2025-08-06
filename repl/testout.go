package repl

import (
	"fmt"
	"unicode/utf8"
)

type TestInput struct {
	buffer   []string
	readHead int
}

func NewTestInput(lines []string) *TestInput {
	return &TestInput{buffer: lines, readHead: 0}
}

func (t *TestInput) Read(p []byte) (n int, err error) {
	if len(t.buffer) == 0 || len(t.buffer) <= t.readHead {
		return 0, fmt.Errorf("Test input has no more lines to read")
	}
	line := fmt.Sprintf("%s\n", t.buffer[t.readHead])
	copy(p, line)
	t.readHead++
	size := 0
	for _, c := range line {
		size += utf8.RuneLen(c)
	}
	return size, nil
}

type TestPrinter struct {
	buffer       []string
	expectations []string
}

func NewTestPrinterExpectations(expectedLines []string) *TestPrinter {
	return &TestPrinter{expectations: expectedLines}
}

func (t *TestPrinter) Println(s string) {
	t.Print(fmt.Sprintf("%s\n", s))
}

func (t *TestPrinter) Print(s string) {
	t.buffer = append(t.buffer, s)
}

func (t TestPrinter) IsValid() error {
	if len(t.buffer) != len(t.expectations) {
		return fmt.Errorf("Test printer recieved %d lines, but wanted %d. Have %v, want %v.", len(t.buffer), len(t.expectations), t.buffer, t.expectations)
	}
	for i, l := range t.buffer {
		if t.expectations[i] != l {
			return fmt.Errorf("Test printer does not match expectations. Have %v, want %v.", t.buffer, t.expectations)
		}
	}
	return nil
}
