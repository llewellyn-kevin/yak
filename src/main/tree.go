package main

type Syntaxtree struct {
  Root *Treenode
}

type Treenode struct {
  Value Token
  Nodes []*Treenode
}

func Newtree(rootval Token) *Syntaxtree {
  tree := new(Syntaxtree)
  tree.Root = new(Treenode)
  tree.Root.Value = rootval
  return tree
}

func (s Syntaxtree) String() string {
  return s.Root.String()
}

func (t *Treenode) Addnode(newval Token) {
  var newnode *Treenode = new(Treenode)
  newnode.Value = newval
  t.Nodes = append(t.Nodes, newnode)
}

func (t *Treenode) Clearnodes() {
  t.Nodes = nil
}

const TAB string = "    "

func (t Treenode) String() (res string) {
  res += t.Value.String()
  if len(t.Nodes) > 0 {
    res += " ("
    for _, value := range t.Nodes {
      res += "\n"
      res += TAB + value.nestedstring(1)
    }
    res += "\n"
    res += ")"
  }
  return
}

func (t Treenode) nestedstring(level int) (res string) {
  res += t.Value.String()
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
