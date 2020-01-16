# yak
A large domesticated ox with shaggy hair.

As a programming language, it functions as an interpreted, 
stack-based, postfix notation language created in go. 

# How to Run
If you compile the main package in go, the result will be
an executable that takes one argument-the name of the file
you wish to interpret. So if you wanted to run the 
recursion program in this repo you would run:

```
$ go build main 
$ ./main recursion.yak
```

The interpreter is still in active development so this will 
be replaced by a more permanent solution in the long run.

The output of the program will be the stack that results 
from running the code. 

# How to Write
Every value the interpreter finds in a program gets placed 
on the stack. Each operation interacts with the values on 
the stack, in most cases popping some off and, and always
placing the result on the top. 

For example: the `+` binary operator, takes the two top
values from the stack, adds them, and pushes the result to 
the stack. 

So a program that reads: 
```
1| 5 10
```
will result in a stack:
```
0: 10
1: 5
``` 

But a program that reads:
```
1| 5 10 +
``` 
will result in a stack:
```
0: 15
```

The lexer seperates all tokens by whitespace. This means 
all inputs and operations must be seperated by spaces, 
linebreaks, or tabs. Which is used is irrelavent. That 
means:

```
1| 1 
2| 2
3| +
4| 3
5| *
```
and 
```
1| 1 2
2| + 3      *
``` 
and even
```
1| 1 2 + 3 *
```
are all valid yak, and do the same thing. But:
```
1| 1 2+ 3* 
```
is invalid. 


## Operators
The binary operators all take two values and return one, 
and are: `+, -, *, /, %, ==`

`+` adds the two values and returns the result 

`-` subtracts the top value from the second value and 
returns the result

`*` multiplies the two values and returns the result 

`/` divides the second value by the top value and returns 
the result 

`%` divides the second value by the top value and returns 
the remainder (expects integers)

`==` compares the two values and returns 1 if they are the 
same, otherwise 0 

The unary operators take one value and return one, the only 
unary operator currently is: `.` 

`.` The duplication operator. This is necessary because 
in most cases interacting with the stack observes values, 
and in most cases, observing the thing also destroys the 
thing. This operator copies the top value from the stack, 
and adds the duplicate to the stack. 

## Control Statements
There are currently 3 control statements: a function 
definition, an if statement, and an if not statement. 

After each control statement, yak requires a block of code. 
A block is defined by an open bracket, `{`, a set of 
instructions, and then a close bracket, `}`. 

### Functions
A function definition takes the form: `n#identifier`, where
n is the number of arguments, and the identifier is the name 
of the function. 

A function is called simply by placing the function name 
in a block of code when there are at least `n` arguments on 
the stack. 

When a function is called the interpreter creates a new 
stack outside the main one named after the function. 
it then pops `n` values off the main stack, and pushes them 
onto the new function stack. After running the code in the 
function block, the top value from the function stack will 
returned to the main stack.

Example: 
```
1| 1#increment {
2|   1 +
3| }
```
This increment function takes one argument, adds one, 
and returns the result. 

### Conditionals
An if statement is simply the ternary operator from 
other languages: `?`. When the interpreter sees the 
conditional operator, it pops off the top of the stack 
and checks its value. If it is 1, it executes the following
code block. If it is not 1, it does not. 

Example: 
```
1| 10
2| 5 5 == ? { 1 + }
```
This will push 10, 5, and 5 onto the stack. Then the `==`
will compare the top two values, see they are the same and 
replace them with `1`. The `?` will pop the `1` off the 
stack and execute the code block. This code block pushes 
a `1` onto the stack and adds it to the remaining value,
ultimately returning 11.

Example 2:
```
1| 10
2| 6 5 == ? { 1 + }
```
Ultimately this is the same example, but `6 5 ==`
returns `0` instead of `1`. This means the block after 
`?` never executes and 10 is ultimately returned.

The if not statement is simply a `!`. It functions the same 
as the conditional operator, except the block of code 
executes if the value observed on the stack is not 1.

Example 3:
```
1| 10
2| 6 5 == ! { 1 + }
```
In this example, the value 0 will be checked, which is 
not 1, so the block will execute. This would also work if 
10 was checked. 

Example 4: 
```
1| 10 . ! { 1 + }
```
Here we duplicate the 10 and check its value. Because it 
is not 1, it is false, and the block of code executes. 
The duplicated 10 was popped of the stack by the `!`,
but we preserved the origional 10 because of the `.`. 
This means we get `10 1 +`, or the expected value: 11.
