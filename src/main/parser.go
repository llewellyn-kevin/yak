package main

import "strings"

func Parse(lexer *Lexer, root *Treenode, endsymbol string) {
  var curtoken Token
  toadd := make(map[string]interface{})
  for ; lexer.Gettoken().Symbol != endsymbol; lexer.Next() {
    curtoken = *lexer.Gettoken()
    toadd = make(map[string]interface{})

    switch(curtoken.Symbol) {
      case "INTEGER":
        toadd["PUSH_INT"] = curtoken.Value
      case "FLOAT":
        toadd["PUSH_FLOAT"] = curtoken.Value
      case "OPERATOR":
        toadd["OPERATION"] = curtoken.Value
      case "FUNC_HEADER":
        toadd["FUNC_HEADER"] = curtoken.Value
        fields := strings.SplitN(curtoken.Value, "#", 2)
        toadd["NUM_ARGS"] = fields[0]
        toadd["FUNC_IDENTIFIER"] = fields[1]

        lexer.Next()
        curtoken = *lexer.Gettoken()
        if curtoken.Symbol != "BLOCK_OPEN" {
          toadd["MISSING_TOKEN"] = "{"
        } else {
          root.Addnode(toadd)
          toadd = make(map[string]interface{})
          lexer.Next()
          Parse(lexer, root.Lastnode(), "BLOCK_CLOSE")
        }
      case "IDENTIFIER":
        toadd["ILLEGAL_TOKEN"] = curtoken.Value
      default:
        toadd["UNRECOGNIZED_TOKEN"] = curtoken.Value
    }
    root.Addnode(toadd)
  }
}

/*func nestedparse(lexer *Lexer, root *Treenode, endsymbol string) {
  car curtoken Token
  toadd := make(map[string]interface{})
  for ; lexer.Gettoken().Value != endsymbol; lexer.Next() {
    curtoken = *lexer.Gettoken()
    toadd = make(map[string]interface{})

    switch(curtoken.Symbol) {
      case
    }
    root.Addnode(toadd)
  }
}*/
