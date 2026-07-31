package eval_test

import (
	"errors"
	"llewellyn-kevin/yak/ast"
	"llewellyn-kevin/yak/eval"
	"llewellyn-kevin/yak/lexer"
	"llewellyn-kevin/yak/parser"
	"llewellyn-kevin/yak/parser/rd"
	"strings"
	"testing"
)

type basicProgramCase struct {
	Input    string
	Expected eval.EvalState
}

func TestBasicProgramState(t *testing.T) {
	input := `%foo 2.3 1 true false`

	expected := makeProgram(
		[]eval.Value{symVal("foo"), floatVal(2.3), intVal(1), boolVal(true), boolVal(false)},
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

func TestUnaryOperator(t *testing.T) {
	input := `1++ -- ++ 2 ++++++ 3 ---- 3.2++`

	expected := makeProgram(
		[]eval.Value{intVal(2), intVal(5), intVal(1), floatVal(4.2)},
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

func TestBasicBinaryOperators(t *testing.T) {
	cases := map[string]basicProgramCase{
		"addition": {
			Input: "1 2 + 4.0 2.4 + 1 2.5 + 4 5.0 + 4.0 6 + -8 2 +",
			Expected: makeProgram(
				[]eval.Value{intVal(3), floatVal(6.4), floatVal(3.5), floatVal(9), floatVal(10), intVal(-6)},
				map[string][]eval.Value{},
			),
		},
		"subtraction": {
			Input: "2 1 - 4.0 2.4 - 5 2.5 - 4 5.0 - 4.0 6 - -3.2 3 -",
			Expected: makeProgram(
				[]eval.Value{intVal(1), floatVal(1.6), floatVal(2.5), floatVal(-1), floatVal(-2), floatVal(-6.2)},
				map[string][]eval.Value{},
			),
		},
		"multiplication": {
			Input: "2 1 * 4.0 2.4 * 5 2.5 * 4 5.0 * 4.0 6 * 3 -2 *",
			Expected: makeProgram(
				[]eval.Value{intVal(2), floatVal(9.6), floatVal(12.5), floatVal(20), floatVal(24), intVal(-6)},
				map[string][]eval.Value{},
			),
		},
		"division": {
			Input: "10 5 / 10.0 2.0 / 5 2 / 5 2.0 / 5.0 2 / 10 -2 /",
			Expected: makeProgram(
				[]eval.Value{intVal(2), floatVal(5), intVal(2), floatVal(2.5), floatVal(2.5), intVal(-5)},
				map[string][]eval.Value{},
			),
		},
		"modulo": {
			Input: "10 3 % 9.3 3 % 10.5 3 % 11 -2 % -11 2 %",
			Expected: makeProgram(
				[]eval.Value{intVal(1), floatVal(.3), floatVal(1.5), intVal(1), intVal(-1)},
				map[string][]eval.Value{},
			),
		},
		"exponent": {
			Input: "",
			Expected: makeProgram(
				[]eval.Value{},
				map[string][]eval.Value{},
			),
		},
		"binary": {
			Input: "12 10 & 12 10 | 12 10 ^ 5 2 << 20 2 >>",
			Expected: makeProgram(
				[]eval.Value{intVal(8), intVal(14), intVal(6), intVal(20), intVal(5)},
				map[string][]eval.Value{},
			),
		},
		"comparison": {
			Input: "4 5 < 4 5 > 5 5 <= 6 5 >= 0 -1 <= 1 1 = 1 2 = true 1 = 1.2 1.2 = 1.0 1 =",
			Expected: makeProgram(
				[]eval.Value{boolVal(true), boolVal(false), boolVal(true), boolVal(true), boolVal(false), boolVal(true), boolVal(false), boolVal(false), boolVal(true), boolVal(false)},
				map[string][]eval.Value{},
			),
		},
	}

	for caseName, data := range cases {
		output, err := evalProgram(data.Input)
		if err != nil {
			t.Errorf("Failed test case %s", caseName)
			t.Errorf("Did not expect any errors. Got %s", err)
			continue
		}

		if !programStatesAreEqual(*output, data.Expected) {
			t.Errorf("Failed test case %s", caseName)
			failOnDivergentPrograms(t, *output, data.Expected)
		}
	}
}

func TestBlocks(t *testing.T) {
	cases := map[string]basicProgramCase{
		"executes a block on main": {
			Input: "1 { 2 } 3",
			Expected: makeProgram(
				[]eval.Value{intVal(1), intVal(2), intVal(3)},
				map[string][]eval.Value{},
			),
		},
		"executes conditional blocks only when true": {
			Input: "true if { 1 } false if { 2 }",
			Expected: makeProgram(
				[]eval.Value{intVal(1)},
				map[string][]eval.Value{},
			),
		},
		"executes else statement on falsey value": {
			Input: "true if { 1 } false if { 2 } else { 3 }",
			Expected: makeProgram(
				[]eval.Value{intVal(1), intVal(3)},
				map[string][]eval.Value{},
			),
		},
		"evaluates any value's truthiness": {
			Input: "1 if { 1 } 0 if { 0 } 4.2 if { 4.2 } 0.0 if { 0.0 } %foo if { %foo } else { %bar }",
			Expected: makeProgram(
				[]eval.Value{intVal(1), floatVal(4.2), symVal("foo")},
				map[string][]eval.Value{},
			),
		},
	}

	for caseName, data := range cases {
		output, err := evalProgram(data.Input)
		if err != nil {
			t.Errorf("Failed test case %s", caseName)
			t.Errorf("Did not expect any errors. Got %s", err)
			continue
		}

		if !programStatesAreEqual(*output, data.Expected) {
			t.Errorf("Failed test case %s", caseName)
			failOnDivergentPrograms(t, *output, data.Expected)
		}
	}
}

func TestAssigment(t *testing.T) {
	cases := map[string]basicProgramCase{
		"basic assignment": {
			Input: "1 ->foo",
			Expected: makeProgram(
				[]eval.Value{},
				map[string][]eval.Value{
					"foo": {intVal(1)},
				},
			),
		},
		"basic n assignment": {
			Input: "1 2 3 3=>foo",
			Expected: makeProgram(
				[]eval.Value{},
				map[string][]eval.Value{
					"foo": {intVal(1), intVal(2), intVal(3)},
				},
			),
		},
		"consecutive assignment": {
			Input: "42 1 4 2 3 2=>foo ->fizz->buzz->fizz",
			Expected: makeProgram(
				[]eval.Value{},
				map[string][]eval.Value{
					"foo":  {intVal(2), intVal(3)},
					"fizz": {intVal(4), intVal(42)},
					"buzz": {intVal(1)},
				},
			),
		},
		"it pulls values off named stacks": {
			Input: "42 1 4 2 3 2=>foo ->fizz->buzz->fizz foo fizz",
			Expected: makeProgram(
				[]eval.Value{intVal(3), intVal(42)},
				map[string][]eval.Value{
					"foo":  {intVal(2)},
					"fizz": {intVal(4)},
					"buzz": {intVal(1)},
				},
			),
		},
	}

	for caseName, data := range cases {
		output, err := evalProgram(data.Input)
		if err != nil {
			t.Errorf("Failed test case %s", caseName)
			t.Errorf("Did not expect any errors. Got %s", err)
			continue
		}

		if !programStatesAreEqual(*output, data.Expected) {
			t.Errorf("Failed test case %s", caseName)
			failOnDivergentPrograms(t, *output, data.Expected)
		}
	}
}

func TestStackOperations(t *testing.T) {
	cases := map[string]basicProgramCase{
		"duplicate single value": {
			Input: "5 .",
			Expected: makeProgram(
				[]eval.Value{intVal(5), intVal(5)},
				map[string][]eval.Value{},
			),
		},
		"duplicate top of stack": {
			Input: "1 2 3 .",
			Expected: makeProgram(
				[]eval.Value{intVal(1), intVal(2), intVal(3), intVal(3)},
				map[string][]eval.Value{},
			),
		},
		"swap two values": {
			Input: "1 2 <>",
			Expected: makeProgram(
				[]eval.Value{intVal(2), intVal(1)},
				map[string][]eval.Value{},
			),
		},
		"swap top of stack": {
			Input: "1 2 3 <>",
			Expected: makeProgram(
				[]eval.Value{intVal(1), intVal(3), intVal(2)},
				map[string][]eval.Value{},
			),
		},
		"swap mixed types": {
			Input: "%foo 2.3 <>",
			Expected: makeProgram(
				[]eval.Value{floatVal(2.3), symVal("foo")},
				map[string][]eval.Value{},
			),
		},
		"combine swaps": {
			Input: "1 2 <> 3 4 <>",
			Expected: makeProgram(
				[]eval.Value{intVal(2), intVal(1), intVal(4), intVal(3)},
				map[string][]eval.Value{},
			),
		},
	}

	for caseName, data := range cases {
		output, err := evalProgram(data.Input)
		if err != nil {
			t.Errorf("Failed test case %s", caseName)
			t.Errorf("Did not expect any errors. Got %s", err)
			continue
		}

		if !programStatesAreEqual(*output, data.Expected) {
			t.Errorf("Failed test case %s", caseName)
			failOnDivergentPrograms(t, *output, data.Expected)
		}
	}
}

func TestStackOperationErrors(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{"empty stack for duplicate", `.`},
		{"empty stack for swap", `<>`},
		{"too few items for swap", `5 <>`},
	}

	for _, tc := range cases {
		_, errs := evalProgram(tc.input)
		if len(errs) == 0 {
			t.Errorf("%s: expected error, got none", tc.name)
			continue
		}
		var rt eval.RuntimeError
		if !errors.As(errs[0], &rt) {
			t.Errorf("%s: got %T, want RuntimeError", tc.name, errs[0])
		}
	}
}

func TestRuntimeErrors(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{"empty stack for unary operator", `++`},
		{"empty stack for conditional", `if { 1 }`},
		{"empty stack for binary operator", `1 2 + +`},
		{"empty stack for bitwise operator", `1 2 & &`},
	}

	for _, tc := range cases {
		_, errs := evalProgram(tc.input)
		if len(errs) == 0 {
			t.Errorf("%s: expected error, got none", tc.name)
			continue
		}
		var rt eval.RuntimeError
		if !errors.As(errs[0], &rt) {
			t.Errorf("%s: got %T, want RuntimeError", tc.name, errs[0])
		}
	}
}

func TestTypeErrors(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{"float bitwise", `1.2 2.3 &`},
		{"symbol + int", `%foo 2 +`},
		{"bool + bool", `true true +`},
		{"bool unary", `true ++`},
	}

	for _, tc := range cases {
		_, errs := evalProgram(tc.input)
		if len(errs) == 0 {
			t.Errorf("%s: expected error, got none", tc.name)
			continue
		}
		var te eval.TypeError
		if !errors.As(errs[0], &te) {
			t.Errorf("%s: got %T, want TypeError", tc.name, errs[0])
		}
	}
}

func TestFunctions(t *testing.T) {
	cases := []basicProgramCase{
		{
			Input: "2#adder#1 { + } 2 3 6 adder adder",
			Expected: makeProgram(
				[]eval.Value{intVal(11)},
				map[string][]eval.Value{},
			),
		},
		{
			Input: "1#addFive#1 { 5 + } 3 addFive",
			Expected: makeProgram(
				[]eval.Value{intVal(8)},
				map[string][]eval.Value{},
			),
		},
		{
			Input: "1#conditionalTest#1 {2 % if {1} else {2}} 4 conditionalTest 9 conditionalTest",
			Expected: makeProgram(
				[]eval.Value{intVal(2), intVal(1)},
				map[string][]eval.Value{},
			),
		},
		{
			Input: "1#double#1 {2 *} 1#halve#1 {2.0 /} 4 double 8 halve",
			Expected: makeProgram(
				[]eval.Value{intVal(8), floatVal(4)},
				map[string][]eval.Value{},
			),
		},
		{
			Input: "1#recursion#1 { 1 - . 0 <= if { } else { recursion } } 5 recursion",
			Expected: makeProgram(
				[]eval.Value{intVal(0)},
				map[string][]eval.Value{},
			),
		},
	}

	for _, tc := range cases {
		output, errs := evalProgram(tc.Input)
		if len(errs) > 0 {
			t.Errorf("Unexpected errors: %v", errs)
			continue
		}

		if !programStatesAreEqual(*output, tc.Expected) {
			failOnDivergentPrograms(t, *output, tc.Expected)
		}
	}
}

func TestFunctionErrors(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{"function not found", `1#foo#1 { + } bar`},
		{"not enough args", `2#foo#2 { + } 1 foo`},
		{"not enough returns", `1#foo#3 { 1 } 0 foo`},
	}

	for _, tc := range cases {
		_, errs := evalProgram(tc.input)
		if len(errs) == 0 {
			t.Errorf("%s: expected error, got none", tc.name)
			continue
		}
		var rt eval.RuntimeError
		if !errors.As(errs[0], &rt) {
			t.Errorf("%s: got %T, want RuntimeError", tc.name, errs[0])
		}
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
func boolVal(v bool) eval.BooleanValue    { return eval.BooleanValue{Value: v} }

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
