package main

import "fmt"

type Stack struct {
  Name string
  Top *Stacknode
}

type Stacknode struct {
  Value interface{}
  Prev *Stacknode
}

func (s Stack) Peek() interface{} {
  return s.Top.Value
}

func (s *Stack) Pop() (top interface{}) {
  top = s.Top.Value
  s.Top = s.Top.Prev
  return
}

func (s *Stack) Put(value interface{}) {
  s.Top = &Stacknode{value, s.Top}
}

func (s Stack) String() string {
  var ret string = s.Name + " {"
  var i int = 0
  for t := s.Top; t != nil; t = t.Prev {
    ret += fmt.Sprintf("\n  %d: %v", i, t.Value)
    i++
  }
  ret += "\n}"
  return ret
}
