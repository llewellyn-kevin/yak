package eval_test

import (
	"llewellyn-kevin/yak/ast"
	"llewellyn-kevin/yak/eval"
	"llewellyn-kevin/yak/lexer"
	"llewellyn-kevin/yak/parser"
	"llewellyn-kevin/yak/parser/rd"
	"strings"
	"testing"
)

func TestBasicProgramState(t *testing.T) {
	input := `%foo 2.3 1`

	expected := makeProgram(
		[]eval.Value{symVal("foo"), floatVal(2.3), intVal(1)},
		map[string][]eval.Value{},
	)

	output, err := evalProgram(input)
	if err != nil {
		t.Errorf("Did not expect any errors. Got %s", err)
	}

	if !programStatesAreEqual(*output, expected) {
		failOnDivergentPrograms(t, *output, expected)
	}
}

func evalProgram(input string) (*eval.EvalState, []error) {
	program, err := parseProgram(input)
	if err != nil {
		return nil, []error{err}
	}
	return eval.Eval(program)
}

func parseProgram(input string) (*ast.Program, error) {
	l := lexer.NewLexer(strings.NewReader(input))
	p, err := parser.NewParserFactory().Use(rd.RECURSIVE_DESCENT_STRATEGY).Get(l)
	if err != nil {
		return nil, err
	}
	return p.Parse(), nil
}

func intVal(v int) eval.IntValue          { return eval.IntValue{Value: v} }
func floatVal(v float64) eval.FloatValue  { return eval.FloatValue{Value: v} }
func stringVal(v string) eval.StringValue { return eval.StringValue{Value: v} }
func symVal(v string) eval.SymbolValue    { return eval.SymbolValue{Value: v} }

func makeProgram(main []eval.Value, named map[string][]eval.Value) (s eval.EvalState) {
	s = *eval.NewEvalState()
	for _, v := range main {
		s.MainStack.Push(v)
	}

	for name, vals := range named {
		newStack := &eval.Stack{}
		for _, v := range vals {
			newStack.Push(v)
		}
		s.NamedStacks[name] = newStack
	}
	return
}

func stacksAreEqual(a, b eval.Stack) bool {
	return a.String() == b.String()
}

func programStatesAreEqual(a, b eval.EvalState) bool {
	if !stacksAreEqual(*a.MainStack, *b.MainStack) {
		return false
	}

	for name, stack := range a.NamedStacks {
		var other *eval.Stack
		var ok bool
		if other, ok = b.NamedStacks[name]; !ok {
			return false
		}
		if !stacksAreEqual(*stack, *other) {
			return false
		}
	}

	return true
}

func failOnDivergentPrograms(t *testing.T, have, want eval.EvalState) {
	t.Errorf("Actual program state does not match expectation. Have %s, Want %s", have, want)
}
