# Xin language specification

## Introduction

Xin is a purely functional programming language inspired by lisp syntax and CSP ideas about concurrency and data sharing. Xin aspires to be an expressive, concise language built on a small number of simple elements that work well together.

Xin supports proper tail calls and uses tail recursion as the only looping primitive, but provides common vocabulary like `while` and `for` / `range` as standard library forms.

## Syntax

Xin's grammar is simple. Here is the complete BNF grammar for Xin.

```
<name> ::= [A-Za-z0-9.:?!-]*
<number> ::= (0-9)+.(0-9)+ | 0x(0-9a-fA-F)+
<string> ::= "'" (.*) "'"

<atom> ::= <name> | <number> | <string>

<form> ::= "(" <atom> ")"

<program> ::= <form>*
```

### Wildcard identifier (?)

The special-case identifier `?` does not reference anything -- to deference `?` is a compile-time error.

`?` is a wildcard identifier, and can be used in two cases.

1. In a definition of a form that takes several arguments, where we want to ignore certain arguments.
2. In a match-case, where we want to match against a generic value or structure.

### Scope

Xin names are lexically scoped. There are three ways that a new scope can be created in a program.

- A new program file creates its own scope.
- When a new form is defined, the body of the form creates a new lexical scope.
- When one or more expressions are included in a do-form, each expression creates its own lexical scope.

The bind form `:` links a value to a name in the current lexical scope.

A bind followed by a raw name creates a static reference to a value.

```
; take the value `val` and assign it to the name `name`
(: name val)
```

A bind followed by a parenthesized expression defines a new form.

```
; define a form `func` with arguments `a, b, c` and define it to
; be equivalent to the expression `body-forms`
(: (func a b c)
   body-forms)

; note that below is legal, and defines a function that returns
; a constant every time, in this case the int 42.
(: (constant-func)
   42)
```

## Types and values

Xin has these types:

- `string`: Unicode string tha can be interpreted as a byte array / blob
- `int`: 64-bit integer
- `frac`: 64-bit floating point
- `form`: bound expression, i.e. a function
- `list`: a heterogeneous linked list of values
- `stream`: a sink / source stream of values for I/O. Streams and stream operations are no pure.

## Evaluation

### Special forms

Xin has 4 special forms.

- `:`: define a new name in the current lexical scope and set it to reference a given value.
- `set`: find an existing name nearest to the current lexical scope and set it to reference a given value. `set` allows mutating nonlocal values. This may be removed in a later version of Xin, if we don't use it or if it's awkward. The semantics of `set` follow that of `:`, as defined in the _Scope_ section above.
- `do`: concurrently evaluate multiple following expressions.
- `match`: a switch-case.

### Concurrency and streams

Streams are the primitive for constructing concurrent programs and doing I/O in Xin. Streams are blocking, synchronous sinks and sources of values.

The source operator `->` pops one value out of the stream. If the stream has no value in the queue, this form blocks.

The sink operator `<-` pushes one value into the stream. If there are expressions blocking on the stream, it un-blocks the first-queued expression.

For example, the `os::out` stream represents the standard out file of a process. Running

```
(<- os::out 'hello world')
```

will print "hello world" to standard out.

Streams can be used to read and write to files, network sockets, or to arbitrary other programs and sinks/sources of values, like OS signals. When writing concurrent programs, blocking stream IO is also the primary synchronization primitive in Xin.

## Packages and imports

All Xin code exists in _packages_. Xin includes by default a set of standard packages like `os`, `math`, `strings`, and `streams`, but Xin programmers can also define their own packages by writing and referencing files with the import name.
