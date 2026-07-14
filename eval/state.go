package eval

import (
	"errors"
	"fmt"
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
		return nil, errors.New("cannot pop off an empty stack")
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
		return nil, fmt.Errorf("cannot pop '%d' items off a stack with only length '%d'", n, s.length)
	}
	vals := make([]Value, n)
	var i uint16
	for i = 0; i < n; i++ {
		var err error
		vals[i], err = s.Pop()
		if err != nil {
			return nil, errors.New("internal error with stack")
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
	MainStack   *Stack
	NamedStacks map[string]*Stack
	activeStack string
}

func NewEvalState() *EvalState {
	return &EvalState{
		MainStack:   &Stack{},
		NamedStacks: make(map[string]*Stack),
	}
}

func (e EvalState) ActiveStack() (*Stack, error) {
	if e.activeStack == "" {
		return e.MainStack, nil
	}

	if active, ok := e.NamedStacks[e.activeStack]; ok {
		return active, nil
	}

	return nil, fmt.Errorf("Internal Error: The current active stack is set as '%s', but there is no stack with that name.", e.activeStack)
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
