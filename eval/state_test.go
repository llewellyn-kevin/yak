package eval_test

import (
	"errors"
	"llewellyn-kevin/yak/eval"
	"testing"
)

func TestStackPushPop(t *testing.T) {
	s := &eval.Stack{}
	v := intVal(42)
	s.Push(v)
	if s.Size() != 1 {
		t.Fatalf("expected size 1, got %d", s.Size())
	}
	popped, err := s.Pop()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if popped != v {
		t.Errorf("expected %v, got %v", v, popped)
	}
	if s.Size() != 0 {
		t.Errorf("expected size 0 after pop, got %d", s.Size())
	}
}

func TestStackPopEmpty(t *testing.T) {
	s := &eval.Stack{}
	_, err := s.Pop()
	var rt eval.RuntimeError
	if !errors.As(err, &rt) {
		t.Errorf("expected RuntimeError, got %T", err)
	}
}

func TestStackPopN(t *testing.T) {
	s := &eval.Stack{}
	a, b := intVal(1), intVal(2)
	s.Push(a)
	s.Push(b)
	vals, err := s.PopN(2)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if len(vals) != 2 {
		t.Fatalf("expected 2 values, got %d", len(vals))
	}
	if vals[0] != b {
		t.Errorf("expected top of stack %v first, got %v", b, vals[0])
	}
	if vals[1] != a {
		t.Errorf("expected second value %v, got %v", a, vals[1])
	}
}

func TestStackPopNTooMany(t *testing.T) {
	s := &eval.Stack{}
	s.Push(intVal(1))
	_, err := s.PopN(2)
	var rt eval.RuntimeError
	if !errors.As(err, &rt) {
		t.Errorf("expected RuntimeError, got %T", err)
	}
}

func TestStackPeek(t *testing.T) {
	s := &eval.Stack{}
	if v := s.Peek(); v != nil {
		t.Errorf("expected nil on empty stack, got %v", v)
	}
	v := intVal(99)
	s.Push(v)
	if s.Peek() != v {
		t.Errorf("expected %v, got %v", v, s.Peek())
	}
	if s.Size() != 1 {
		t.Errorf("expected size 1 after peek (no removal), got %d", s.Size())
	}
}

func TestStackLIFOOrder(t *testing.T) {
	s := &eval.Stack{}
	vals := []eval.Value{intVal(1), intVal(2), intVal(3)}
	for _, v := range vals {
		s.Push(v)
	}
	for i := len(vals) - 1; i >= 0; i-- {
		popped, err := s.Pop()
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if popped != vals[i] {
			t.Errorf("expected %v, got %v", vals[i], popped)
		}
	}
}

func TestStackSize(t *testing.T) {
	s := &eval.Stack{}
	checkSize := func(want uint16) {
		t.Helper()
		if got := s.Size(); got != want {
			t.Errorf("expected size %d, got %d", want, got)
		}
	}
	checkSize(0)
	s.Push(intVal(1))
	checkSize(1)
	s.Push(intVal(2))
	checkSize(2)
	s.Pop()
	checkSize(1)
	s.Pop()
	checkSize(0)
}

func TestStackString(t *testing.T) {
	s := &eval.Stack{}
	if s.String() != "(count: 0)" {
		t.Errorf("expected empty string, got %s", s.String())
	}
	s.Push(intVal(7))
	s.Push(intVal(42))
	got := s.String()
	want := "(count: 2) 42 7"
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}
