package main

import (
  "reflect"
  "fmt"
)

var funcnodes []*Treenode

func Interpret(root *Treenode, stack *Stack) *Stack {
  for _, node := range root.Nodes {
    for k, v := range node.Value {
      switch k {
      case "PUSH_INT":
        stack.Put(v)
      case "PUSH_FLOAT":
        stack.Put(v)
      case "OPERATION":
        switch v {
        case "+":
          a, erra := getFloat(stack.Pop())
          b, errb := getFloat(stack.Pop())

          if erra == nil && errb == nil {
            stack.Put(a + b)
          }
        case "-":
          a, erra := getFloat(stack.Pop())
          b, errb := getFloat(stack.Pop())

          if erra == nil && errb == nil {
            stack.Put(b - a)
          }
        case "*":
          a, erra := getFloat(stack.Pop())
          b, errb := getFloat(stack.Pop())

          if erra == nil && errb == nil {
            stack.Put(a * b)
          }
        case "/":
          a, erra := getFloat(stack.Pop())
          b, errb := getFloat(stack.Pop())

          if erra == nil && errb == nil {
            stack.Put(b / a)
          }
        case "%":
          a, erra := getFloat(stack.Pop())
          b, errb := getFloat(stack.Pop())
          c := int(a)
          d := int(b)

          if erra == nil && errb == nil {
            stack.Put(d % c)
          }
        case "==":
          a, erra := getFloat(stack.Pop())
          b, errb := getFloat(stack.Pop())

          if erra == nil && errb == nil {
            if a == b {
              stack.Put(1)
            } else {
              stack.Put(0)
            }
          }
        case "<":
          a, erra := getFloat(stack.Pop())
          b, errb := getFloat(stack.Pop())

          if erra == nil && errb == nil {
            if a > b {
              stack.Put(1)
            } else {
              stack.Put(0)
            }
          }
        case "<=":
          a, erra := getFloat(stack.Pop())
          b, errb := getFloat(stack.Pop())

          if erra == nil && errb == nil {
            if a >= b {
              stack.Put(1)
            } else {
              stack.Put(0)
            }
          }
        case ">":
          a, erra := getFloat(stack.Pop())
          b, errb := getFloat(stack.Pop())

          if erra == nil && errb == nil {
            if a < b {
              stack.Put(1)
            } else {
              stack.Put(0)
            }
          }
        case ">=":
          a, erra := getFloat(stack.Pop())
          b, errb := getFloat(stack.Pop())

          if erra == nil && errb == nil {
            if a <= b {
              stack.Put(1)
            } else {
              stack.Put(0)
            }
          }
        case ".":
          stack.Put(stack.Peek())
        }
      case "IF":
        a, err := getFloat(stack.Pop())

        if err == nil {
          if a == 1 {
            stack = Interpret(node, stack)
          }
        }
      case "NOT":
        a, err := getFloat(stack.Pop())

        if err == nil {
          if a != 1 {
            stack = Interpret(node, stack)
          }
        }
      case "FUNC_HEADER":
        funcnodes = append(funcnodes, node)
      case "FUNC_IDENTIFIER", "NUM_ARGS":
      case "FUNC_CALL":
        for _, f := range funcnodes {
          if f.Value["FUNC_IDENTIFIER"] == v {
            args, err := getFloat(f.Value["NUM_ARGS"])
            if err == nil {
              funcstack := new(Stack)
              funcstack.Name = "Function Call"
              var i float64
              for i = 0; i < args; i++ {
                funcstack.Put(stack.Pop())
                funcstack = Interpret(f, funcstack)
                stack.Put(funcstack.Pop())
              }
            }
          }
        }
      }
    }
  }

  return stack
}

var floatType = reflect.TypeOf(float64(0))

func getFloat(unk interface{}) (float64, error) {
    v := reflect.ValueOf(unk)
    v = reflect.Indirect(v)
    if !v.Type().ConvertibleTo(floatType) {
        return 0, fmt.Errorf("cannot convert %v to float64", v.Type())
    }
    fv := v.Convert(floatType)
    return fv.Float(), nil
}
