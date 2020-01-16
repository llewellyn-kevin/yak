# yak
A large domesticated ox with shaggy hair.

As a programming language, it functions as an interpreted, 
stack-based, postfix notation language created in go. 

# How to Run 

### Option 1 - Unix 
If you don't have a go workspace set up but would like to: 
[here is a good article detailing go workspaces.](https://medium.com/rungo/working-in-go-workspace-3b0576e0534a)

If you have a go workspace set up with your `$GOPATH`
set and `$GOBIN` in your path then it will be very easy. 
Simply use `get` to grab this repo and `install` to create 
the binary. Then you can just run the program which takes 
the name of the file you wish to parse as an argument.

For example, if you wanted to run one of the example programs
included in this repo:

```bash
$ go get github.com/llewellyn-kevin/yak 
$ go install yak 
$ cp $GOPATH/src/github.com/llewellyn-kevin/yak/examples/recursion.yak .
$ yak recursion.yak
```

### Option 2 - Unix
If you are not a go developer and do not plan to be you 
won't have the whole workspace set up. You will still have 
to have [go installed](https://golang.org/doc/install), after 
you do you will have to manually compile from this repo:

```bash 
$ git clone https://github.com/llewellyn-kevin/yak.git
$ cd yak 
$ go build yak
```

This will create an executable binary called yak you can 
put anywhere you like (or just leave it).

You can use this executable like normal: 
```bash 
$ ./yak examples/recursion.yak
```
And yak will interpret the recursion yak program. 

Or you can add yak to your path to run it from anywhere on 
your system. 
```bash 
$ export PATH=$PATH:$PWD/yak 
$ cd examples 
$ yak recursion.yak
```
You can add the `export` line with `$PWD` replaced with the 
path to the yak binary to `~/.bash_profile` so the yak 
command will work on terminal startup. 

### Option 3 - Windows
Go works in windows. So download the code, compile it, 
and run the executable. I don't use Windows, so I wouldn't 
even know how. But you do, so I hope you can.

--------------------------------------------------------------------------------
The interpreter is still in active development, when it is 
closer to an official release I will add binaries that can
be downloaded directly without having to compile using go.

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
