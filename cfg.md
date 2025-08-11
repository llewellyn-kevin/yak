# Context-Free Grammar for the Yak Language

Root
```
S -> S* S | epsilon 
S* -> Z | F
Z -> Z* Z | epsilon
Z* -> A | C | W | O | L
```

Assignment - Statements
```
A -> assign ident A*
A* -> colon T | epsilon
```

Conditionals - Statements
```
C -> I | I E | N | N E
I -> if B
N -> not B
E -> else B
B -> lbrace Z rbrace
B' -> lbrace ident Z rbrace
```

For Loops - Statement
```
W -> for B | for B' 
```

Function Definitions - Statements
```
F -> A hash ident F* B
F* -> hash A | epsilon
A -> int | lparen A* rparen
A* -> A' A* | A'
A' -> int colon T
T -> Y Y*
Y* -> pipe Y Y* | pipe Y | epsilon
Y -> ident | kbool | kint | kfloat | kstring | ksymbol
```

Operations - Expressions
```
O -> add | sub | mult | div | mod | lt | gt | lteq | gteq | eq | swap | bor | band | bxor |
     dup | inc | dec | lshift | rshift | yakout | yakin | literal_yakout | yakup
```

Literals - Statements
```
L -> true, false, int, float, string, symbol
```
