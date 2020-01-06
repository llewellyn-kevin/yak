package main

import (
  "reflect"
  "fmt"
)

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
