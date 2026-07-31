package eval

import (
	"fmt"
	"llewellyn-kevin/yak/ast"
)

type node struct {
	Value Value
	next  *node
	prev  *node
}

func newNode(v Value) *node {
	return &node{Value: v}
}

type Stack struct {
	head   *node
	length uint16
}

func NewStackFromValues(vals []Value) *Stack {
	s := &Stack{}
	for _, v := range vals {
		s.Push(v)
	}
	return s
}

func (s *Stack) Push(v Value) {
	n := newNode(v)
	if s.head != nil {
		s.head.prev = n
		n.next = s.head
	}
	s.head = n
	s.length += 1
}

func (s *Stack) Pop() (Value, error) {
	if s.length == 0 {
		return nil, RuntimeError{Message: "cannot pop off an empty stack"}
	}

	v := s.head.Value
	s.head = s.head.next
	if s.head != nil {
		s.head.prev = nil
	}
	s.length -= 1
	return v, nil
}

func (s *Stack) PopN(n uint16) ([]Value, error) {
	if s.length < n {
		return nil, RuntimeErrorf("cannot pop '%d' items off a stack with only length '%d'", n, s.length)
	}
	vals := make([]Value, n)
	for i := range n {
		var err error
		vals[i], err = s.Pop()
		if err != nil {
			return nil, InternalError{Message: "internal error with stack"}
		}
	}
	return vals, nil
}

func (s *Stack) Peek() Value {
	if s.length == 0 {
		return nil
	}

	return s.head.Value
}

func (s Stack) Size() uint16 {
	return s.length
}

func (s Stack) String() (str string) {
	str += fmt.Sprintf("(count: %d)", s.length)
	for n := s.head; n != nil; n = n.next {
		str += " " + n.Value.String()
	}
	return
}

type EvalState struct {
	FunctionTable map[string]*ast.FunctionStatement
	MainStack     *Stack
	NamedStacks   map[string]*Stack
	activeStack   string
}

func NewEvalState() *EvalState {
	return &EvalState{
		MainStack:   &Stack{},
		NamedStacks: make(map[string]*Stack),
	}
}

func (e *EvalState) AddFunction(fn *ast.FunctionStatement) {
	if e.FunctionTable == nil {
		e.FunctionTable = make(map[string]*ast.FunctionStatement)
	}
	e.FunctionTable[fn.Name] = fn
}

func (e *EvalState) AddFunctionsFromMap(funcs map[string]*ast.FunctionStatement) {
	for _, fn := range funcs {
		e.AddFunction(fn)
	}
}

func (e EvalState) HasFunction(name string) bool {
	_, ok := e.FunctionTable[name]
	return ok
}

func (e EvalState) executeFunction(name string) []error {
	fn, ok := e.FunctionTable[name]
	if !ok {
		return []error{RuntimeErrorf("could not find function with name '%s'", name)}
	}

	stack, err := e.ActiveStack()
	prevStackName := e.activeStack
	if err != nil {
		return []error{err}
	}

	vals, err := stack.PopN(uint16(fn.Args))
	if err != nil {
		return []error{RuntimeErrorf("function '%s' expects %d argument(s), but there were not enough values on stack '%s'", name, fn.Args, e.stackLabel())}
	}
	fnStack := NewStackFromValues(vals)
	fnStackName := fmt.Sprintf("fn#%s", name)
	e.NamedStacks[fnStackName] = fnStack
	e.activeStack = fnStackName

	defer func() {
		e.activeStack = prevStackName
		e.RemoveStack(fnStackName)
	}()

	errors := e.executeBlock(fn.Body)
	if len(errors) > 0 {
		return errors
	}

	returnVals, err := fnStack.PopN(uint16(fn.Returns))
	if err != nil {
		return []error{RuntimeErrorf("function '%s' did not return enough values, expected %d but got %d", name, fn.Returns, len(returnVals))}
	}

	for _, val := range returnVals {
		stack.Push(val)
	}

	return []error{}
}

func (e EvalState) ActiveStack() (*Stack, error) {
	if e.activeStack == "" {
		return e.MainStack, nil
	}

	if active, ok := e.NamedStacks[e.activeStack]; ok {
		return active, nil
	}

	return nil, InternalErrorf("the current active stack is set as '%s', but there is no stack with that name", e.activeStack)
}

func (e *EvalState) RemoveStack(name string) error {
	if _, ok := e.NamedStacks[name]; !ok {
		return RuntimeErrorf("could not remove stack '%s', it does not exist", name)
	}
	delete(e.NamedStacks, name)
	return nil
}

func (e EvalState) stackLabel() string {
	r := e.activeStack
	if r == "" {
		return "main"
	}
	return r
}

func (e EvalState) NamedStacksString() (str string) {
	for k, v := range e.NamedStacks {
		str += fmt.Sprintf("\t\t[%s] %s\n", k, v.String())
	}
	return
}

func (e EvalState) String() string {
	return "EvalState {\n" +
		"\tMainStack: " + e.MainStack.String() + "\n" +
		"\tNamedStacks: \n" + e.NamedStacksString() +
		"}"
}
