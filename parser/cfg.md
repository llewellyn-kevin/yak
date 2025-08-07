# Context-Free Grammar for the Yak Language

Statements
```
S -> S* S | epsilon 
S* -> Z | F
Z -> Z* Z | epsilon
Z* -> A | C | W | O | L
```

Assignment
```
A -> assign ident A*
A* -> colon T | epsilon
```

Conditionals
```
C -> I | I E | N | N E
I -> if B
N -> not B
E -> else B
B -> lbrace Z rbrace
B' -> lbrace ident Z rbrace
```

For Loops
```
W -> for B | for B' 
```

Function Definitions
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

Operations
```
O -> add | sub | mult | div | mod | lt | gt | lteq | gteq | eq | swap | bor | band | bxor |
     dup | inc | dec | lshift | rshift | yakout | yakin | literal_yakout | yakup
```

Literals
```
L -> true, false, int, float, string, symbol
```
