package main


func Parse(lexer *Lexer, root *Treenode, endsymbol string) {
  var curtoken Token
  toadd := make(map[string]interface{})
  for ; lexer.Gettoken().Value != endsymbol; lexer.Next() {
    curtoken = *lexer.Gettoken()
    toadd = make(map[string]interface{})

    switch(curtoken.Symbol) {
      case "INTEGER":
        toadd["PUSH_INT"] = curtoken.Value
      case "FLOAT":
        toadd["PUSH_FLOAT"] = curtoken.Value
      case "OPERATOR":
        toadd["OPERATION"] = curtoken.Value
      case "IDENTIFIER":
        toadd["ILLEGAL_TOKEN"] = curtoken.Value
      default:
        toadd["UNRECOGNIZED_TOKEN"] = curtoken.Value
    }
    root.Addnode(toadd)
  }
}
