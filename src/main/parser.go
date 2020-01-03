package main


func Parse(lexer *Lexer, root *Treenode, endsymbol string) {
  for ; lexer.Gettoken().Value != endsymbol; lexer.Next() {
    root.Addnode(*lexer.Gettoken())
  }
}
