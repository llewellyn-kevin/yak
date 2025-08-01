# yak
A large domesticated ox with shaggy hair.

As a programming language, it functions as an interpreted, 
stack-based, postfix notation language created in go. 

## Breakdown of ideas

data structures
- list - just a variable, variables are all named stacks
- dict
- symbols - essentially enum values that can be defined in real time, like ruby. Preceded by a colon

control structures
- conditional
- loop - python style, go through stack until empty, have a way to generate a stack of numbers up to value
- function - take args off the stack, and 

Operators
----------
binary operations (v1 - top of stack, v2 next to top value)
- `+` := v2 + v1
- `-` := v2 - v1
- `*` := v2 \* v1
- `/` := v2 / v1
- `%` := v2 % v1
- `<` := return 1 if v2 < v1 else 0
- `>` := return 1 if v2 > v1 else 0
- `>=` := return 1 if v2 >= v1 else 0
- `<=` := return 1 if v2 <= v1 else 0
- `=` := return 1 if v2 == v1 else 0
- `<>` := swap v1 and v2
- `|` := v2 | v1 (OR)
- `&` := v2 & v1 (AND)
- `^` := v2 ^ v1 (XOR)

unary operations
- `.` := duplicate top value
- `++` := increment top value
- `--` := decrement top value
- `<<` := left shift
- `>>` := right shift

Variables
------
`-> i` := put the value on top of the main stack into stack i
`=> i` := put `n` values in stack i, where n is the current top of the stack. The values are copied in current order
`i` := pop the value of stack i into the main stack
`.i` := read the value off of the i stack and put it on the current main stack (duplicate & pop)
`<>i` := swap the two values on top of stack i
`{i /** code */}` := execute a code block with i as the main stack
`yakup` pop a value off of the parent main stack if in a code block
`yakupup` to get up multiple scopes, up can be repeated
`-> {/** code */}` := create a new anonymous block where the value of the arrow is the value of the main function
`i:types` := specify the types that are allowed to be put in stack i

```
5 -> newStack
// yakout throws an error, main stack empty
.newStack yakout
// Prints: 5
newStack yakout
// Prints: 5
// newStack now prints an error for an empty stack

1 2 3 3 => secondStack
{secondStack
	yakout yakout yakout
	// Prints 1\n2\n3
}

9 8 7 3 => {
	yakout yakout yakout
	// Prints 9\n8\n7
}

'foo' 1 2 -> integerStack:int // Puts 2 in integer stack
-> integerStack // Puts 1 in integer stack
// -> integerStack // InvalidValueError trying to put 'foo' in integerStack
```

Supported Types
-------
`int`
`float`
`string`
`bool`
`symbol`

```
1 -> iStack:int
1.5 -> fStack:float
'one' -> sStack:string
:one -> syStack:symbol

1 :one 2 => unionStack:int|symbol
// 'one' -> unionStack // InvalidValueError

// string literals are interpreted as specific symbol ids
:one -> specificSymbolStack:one|two|three // valid
:three -> specificSymbolStack // valid
// :too -> specificSymbolStack // invalid


```

Conditionals
--------
`if { }` := evaluate the code block if the top of the stack is truthy
`not { }` := evaluate the code block if the top of the stack is falsy
`if { } : { }` := evaluate the first block if the top of the stack is truthy, otherwise execute the second code block

```
1 2 = {
	'truthy' yakout
} : {
	'falsey' yakout
}

// Prints falsey
```

For
-----
`for` := when this keyword appears in a scoped block, then when the end of the block is reached, the code returns to for until the main stack for the block is empty

`min max increment range -> var` := put the numbers between min and max with increment inc into stack var

```
1 10 2 range -> i
{i for
	yakout
}

// Prints: 1, 3, 5, 7, 9
```

Key-Value Store
---------
there is a global key value store that is represented as a hash table in the interpreter. by default this key value store has the settings used for input and output with the yak functions. this can be set and read with `set` and `get` functions:

```
// Schema {value} {collection} {key} set

'bar' :new-collection :foo set

:new-collection :foo get
yakout

// Prints: bar
```

the order of the function may seem counterintuitive, but the value should be first so it can just come from the top of the main stack if needed. The key being last allows autocomplete with known keys for the collection

values in the store are invalidated when the current block exits. To prevent this tag it with a global scope (using `setg` function):

```
42 -> foo
{foo
	. :global-collection :foo setg
	:local-collection :foo set
}

:global-collection :foo get yakout
// Prints: 42
// :local-collection :foo get yakout // throws an error
```

set valid inputs for hash key value. Unknown types are treated as symbols, known types are that type.

```
:collection :ui-mode 'light|dark' setopts
// allows :light or :dark

:yakin :file-scanner 'characters|lines|tokens|string' setopts
// allows any string, :characters, :lines, :tokens
```

Functions
---------
definition: `arg_list#func_name#return_list { }`
to invoke: `func_name`
to invoke and store values in named stack: `func_name -> var`
if you want to invoke functions multiple times add one or several `.` to the end of the identifier

```
', ' :yak out-del set

1#add5#1 {
	5 +
}

10 add5 yakout
// Prints: 15

10 add5.. yakout
// Prints: 25

1 add5 -> added
2 add5 -> added
3 add5 -> added
{added for yakout}
// Prints 6, 7, 8

4 5 6 add5.. -> addedTwo
{addedTwo for yakout}
// Prints 9, 10, 11

(2:int)#simpleAdder#(1:int) {
	+
}

3 6 simpleAdder yakout
// Prints: 9
// '1' 2 simpleAdder throws an invalid argument error
```

`arg_list` := `(2:int 3:string)` --- the initialized stack must be 3 strings than 2 ints
    | `arg_list` := `2` --- the initialized stack takes 2 values of any type
	| `arg_list` := (2:int|sting) --- union types, int or string
`return_list` is the same format as arg list, but can also be empty, returning only 1 value


I/O settings
----------
`yakout` pops the current value and prints to desired output
`yakout!` same as above, but don't add delimiters
`yakin` puts the next input value on the stack

Use the setyak function to determine how to read input and where to send output

- settings:
	- `value :yak key set` := takes a setting and value off the stack and changes current input output methods
- valid `:in-type`
	- `:args` := read the next argv
	- `:file` := treat argument as a file, 
- valid `:file-scanner`
	- `:characters`
	- `:tokens`
	- `:lines`
	- `string` := arbitrary string(s) with a delimiter to use for separating tokens (e.g. ',', '\n' for csv)
- valid `:out-type`
	- `:stdout` := prints to stdout
	- `:file` := prints to file
- valid `:out-file`
	- `string` := string for the current filename to write out to
- valid `:out-del`
	- `string` := arbitrary strings, (e.g. ',' for csv)
	- `:new-line` := newline character
	- `:space` := a single space

Bash Commands
--------
IDK yet



CFG
--------

All operators with their name as token
All keywords with their name as token

OP -> `PLUS|MINUS|MULT|DIV|MOD|IF|RSHIFT` etc...

ASSIGN -> `->`
NASSIGN -> `n->`
LEFT_PAREN -> `(`
RIGHT_PAREN -> `)`
LEFT_BRACKET -> `{`
RIGHT_BRACKET -> `}`
HASH -> `#`

NUMBER -> `[0-9]+`
CHARACTER -> `[a-zA-Z]+`

TYPE -> `bool|int|float|string|symbol|SYMBOL_NAME

IDENTIFIER -> `[CHARACTER|NUMBER|_]+`
SYMBOL_ID -> `[CHARACTER|NUMBER|-]+`

SYMBOL_DECLARATION -> `COLONSYMBOL_ID`

HYPHEN -> `-`

SINGLE_QUOTE -> `'`
DOUBLE_QUOTE -> `"`
STRING_LITERAL -> `SINGLE_QUOTE.*SINGLE_QUOTE|DOUBLE_QUOTE.*DOUBLE_QUOTE
INT_LITERAL -> `HYPHENNUM+|NUM+`
FLOAT_LITERAL -> `HYPHENNUM+.NUM+|NUM+.NUM+`
TRUE -> true
FALSE -> false

TYPED_ARG -> `NUM:TYPE`

ARG_LIST -> `NUM|(TYPED_ARG+)`
VARIABLE -> IDENTIFIER
FUNCTION -> `ARG_LIST#IDENTIFIER|ARG_LIST#IDENTIFIER#ARG_LIST [WHITESPACE|LEFT_BRACKET]`
