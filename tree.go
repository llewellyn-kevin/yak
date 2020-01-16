package main

import "fmt"

type Syntaxtree struct {
  Root *Treenode
}

type Treenode struct {
  Value map[string]interface{}
  Nodes []*Treenode
}

func Newtree(rootval map[string]interface{}) *Syntaxtree {
  tree := new(Syntaxtree)
  tree.Root = new(Treenode)
  tree.Root.Value = rootval
  return tree
}

func (s Syntaxtree) String() string {
  return s.Root.String()
}

func (t *Treenode) Addnode(newval map[string]interface{}) {
  var newnode *Treenode = new(Treenode)
  newnode.Value = newval
  t.Nodes = append(t.Nodes, newnode)
}

func (t *Treenode) Clearnodes() {
  t.Nodes = nil
}

func (t Treenode) Lastnode() *Treenode {
  return t.Nodes[len(t.Nodes) - 1]
}

const TAB string = "    "

func (t Treenode) String() (res string) {
  for k, v := range t.Value {
    res += fmt.Sprintf("%s: %v", k, v)
  }
  if len(t.Nodes) > 0 {
    res += " ("
    for _, value := range t.Nodes {
      res += "\n"
      res += value.nestedstring(1)
    }
    res += "\n"
    res += ")"
  }
  return
}

func (t Treenode) nestedstring(level int) (res string) {
  for k, v := range t.Value {
    res += "\n"
    for i := 0; i < level; i++ { res += "   " }
    res += fmt.Sprintf("%s: %v", k, v)
  }
  if len(t.Nodes) > 0 {
    res += " ("
    for _, value := range t.Nodes {
      res += "\n"
      for i := 0; i < level; i++ { res += TAB }
      res += TAB + value.nestedstring(level+1)
    }
    res += "\n"
    for i := 0; i < level; i++ { res += TAB }
    res += ")"
  }
  return

}
