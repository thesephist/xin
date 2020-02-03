# Xin language specification

## Introduction

Xin is a functional programming language inspired by Lisp and CSP. Xin aspires to be an expressive, extensible language built on a small number of simple elements that work well together.

Xin supports proper tail calls and uses tail recursion as the only looping primitive, but provides common vocabulary like `while` and `for` / `range` as standard library forms. Xin uses lazy evaluation.

## Syntax

Xin's grammar is simple. Here is the complete BNF grammar for Xin. Notably, any non-whitespace non-delimiter (`()[]{}`) character is valid in an identifier.

```
<name> ::= [\w\-\?\!\+\*\/\:\>\<\=\%\&\|]*
<number> ::= ("-" | "") (0-9)+.(0-9)+ | 0x(0-9a-fA-F)+
<string> ::= "'" (.*) "'"

<atom> ::= <name> | <number> | <string>

<form> ::= "(" (<atom> | <form>)+ ")"

<program> ::= <form>*
```

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

`form`, `vec`, `map`, and `stream` types are passed and equality-checked by reference, all others are passed and equality checked by value.

## Evaluation

### Special forms

Xin has 4 special forms.

- `:`: define a new name in the current lexical scope and set it to reference a given value
- `if`: an if-else
- `do`: sequentially evaluate multiple following expressions
- `import`: import external files as Xin programs

### Streams

Streams are the primitive for constructing concurrent programs and doing I/O in Xin. Streams are sinks and sources of values that interface with the rest of the host system, or another remote part of the Xin program.

The source operator `->` pops one value out of the stream.

The sink operator `<-` pushes one value into the stream.

For example, the `os::stdout` stream represents the standard out file of a process. Running

```
(<- os::stdout 'hello world\n')
```

will print "hello world" to standard out.

Streams can be used to read and write to files, network sockets, or to arbitrary other programs and sinks/sources of values, like OS signals. Despite the vocabulary for concurrency in Xin, Xin is a single-threaded language with an interpreter lock, and concurrent calls are ordered onto a single execution timeline, like JavaScript and CPython asyncio.

### Lazy evaluation

Evaluation of all expressions and forms in Xin are deferred until they are coerced to resolve to a value by a special form, or by returning a value to the Repl. Expressions are coerced to take real values in the following cases:

1. Before being bound to a name via the bind (`:`) special form
2. As a top-level expression in a `do` special form
3. When interfacing with any runtime API that interacts with the outside world, like `os::stdout`
4. When returned as a value to the Repl.

## Packages and imports

Xin includes by default a set of standard packages like `os`, `math`, `vec`, and `str` (string). Xin programmers can also define their own packages by writing and referencing files with the import name.

A Xin program can import values defined in another Xin program with the `import` form. There are two ways to import.

- `(import path)`: find file described by `path` and make all values defined in that file available under the current global namespace.
- `(import path alias)`: find file described by `path` and make values available under _only_ the alias `alias::` namespace.

```
; makes forms and values defined in file 'core/models/user.xin'
; available under the namespace `user-model::xxx`
(import 'core/models/user' user-model)
```
## Proposals

This proposals section documents ideas about additions and changes to the language that I have considered or others have suggested, but I don't think are well aligned with the current values of the language or are too early in their specifications to be included in the implementation or specification.

### Spread syntax

Spread syntax enables Xin programmers to "collect" or "spread" arbitrary segments of the arguments to a form in its invocation or definition, and refer to it as a vector.

### Composite type literals

Add literal syntax for vectors and maps.

Vector literals are delimited with square brackets:

```
(: my-vec [1 2 3 4 [5 6]])

; equivalent to

(: my-vec
   (vec 1 2 3 4 (vec 5 6)))
```

Map literals are delimited with braces and use the arrow `->`:

```
(: letter-counts
   {'hi' -> 2 'hello' -> 5 'bye' -> 3})

; equivalent to

(: letter-counts
   (do (: m (map))
     (map::set! m 'hi' 2)
     (map::set! m 'hello' 5)
     (map::set! m 'bye' 3)))
```

