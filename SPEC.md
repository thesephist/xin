# Xin language specification

## Introduction

Xin is a purely functional programming language inspired by lisp syntax and CSP ideas about concurrency and data sharing. Xin aspires to be an expressive, concise language built on a small number of simple elements that work well together.

Xin supports proper tail calls and uses tail recursion as the only looping primitive, but provides common vocabulary like `while` and `for` / `range` as standard library forms.

## Syntax

Xin's grammar is simple. Here is the complete BNF grammar for Xin. Notably, any non-whitespace non-delimiter (`()[]{}`) character is valid in an identifier.

```
<name> ::= [\w\-\?\!\+\*\/\:\>\<\=\%\&\|]*
<number> ::= ("-" | "") (0-9)+.(0-9)+ | 0x(0-9a-fA-F)+
<string> ::= "'" (.*) "'"

<atom> ::= <name> | <number> | <string>

<form> ::= "(" <name> <atom>* ")"

<program> ::= <form>*
```

### Wildcard identifier (?)

The special-case identifier `?` does not reference anything -- to deference `?` is a compile-time error.

`?` is a wildcard identifier, and can be used in two cases.

1. In a definition of a form that takes several arguments, where we want to ignore certain arguments.
2. In an equality check or comparison, where we want to compare against a generic value or structure that will always return true.

### Scope

Xin names are lexically scoped. There are two ways that a new scope can be created in a program.

- A new program file creates its own scope.
- When a new form is defined, the body of the form creates a new lexical scope.

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
- `vec`: a heterogeneous list of values
- `map`: a heteroogenous hashmap of values
- `stream`: a sink / source stream of values for I/O. Stream operations are not pure.

## Evaluation

### Special forms

Xin has 3 special forms.

- `:`: define a new name in the current lexical scope and set it to reference a given value
- `if`: an if-else
- `do`: sequentially evaluate multiple following expressions

### Concurrency and streams

Streams are the primitive for constructing concurrent programs and doing I/O in Xin. Streams are blocking, synchronous sinks and sources of values.

The source operator `->` pops one value out of the stream. If the stream has no value in the queue, this form blocks.

The sink operator `<-` pushes one value into the stream. If there are expressions blocking on the stream, it un-blocks the first-queued expression.

For example, the `os::out` stream represents the standard out file of a process. Running

```
(<- os::stdout 'hello world')
```

will print "hello world" to standard out.

Streams can be used to read and write to files, network sockets, or to arbitrary other programs and sinks/sources of values, like OS signals. When writing concurrent programs, blocking stream IO is also the primary synchronization primitive in Xin.

Despite the vocabulary for concurrency in Xin, Xin is a single-threaded language with an interpreter lock, and concurrent calls are ordered onto a single execution timeline.

## Packages and imports

Xin includes by default a set of standard packages like `os`, `math`, `string`, and `stream`. Xin programmers can also define their own packages by writing and referencing files with the import name.

A Xin program can import values defined in another Xin program with the `import` form. There are two ways to import.

- `(import import-specifier)`: find file described by `import-specifier` and make all values defined in that file available under the `import-specifier::` namespace.
- `(import import-specifier alias)`: find file described by `import-specifier` and make values available under _only_ the alias `alias::` namespace.

```
; import forms defined in math standard package
; makes math::sqrt, math::sin, math::cos, etc available
(import math)
; makes forms defined in file 'core/models/user.xin'
; available under the namespace `user-model::xxx`
(import core::models::user user-model)
```
