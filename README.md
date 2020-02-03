# Xin programming language

[![GoDoc](https://godoc.org/github.com/thesephist/xin?status.svg)](https://godoc.org/github.com/thesephist/xin)
[![Build Status](https://travis-ci.com/thesephist/xin.svg?branch=master)](https://travis-ci.com/thesephist/xin)

Xin is a functional programming language inspired by Lisp and CSP. Xin aspires to be an expressive, extensible language built on a small number of simple elements that work well together. You can find a deep dive into the language design [in the spec document](SPEC.md).

Xin is my second toy programming language, after [Ink](https://github.com/thesephist/ink). With Xin, I'm specifically exploring ideas around code as a data structure, lazy evaluation in an interpreter context, and streaming evented I/O.

Here's the fibonacci sequence, written naively in Xin.

```
(: (fib n)
   (if (| (= n 0) (= n 1))
     1
     (+ (fib (- n 1))
        (fib (- n 2)))))

(log (fib 20))
```

Xin supports proper tail calls, so we can write this in a faster (`O(n)`) tail-recursive form.

```
(: (fibh n a b)
   (if (| (= n 0) (= n 1))
     b
     (fibh (- n 1) b (+ a b))))
(: (fib n)
   (fibh n 1 1))

(log (fib 20))
```

## Goals

- Expressive, extensible syntax that's natural to read and suitable for defining DSLs
- Programs that lend themselves to sophisticated data structures and dumb algorithms, rather than vice versa
- Ease of learning
- Great interactive Repl

## Installation and usage

While Xin is under development, the best way to get it is to build it from source. If you have Go installed, clone the repository and run

```
go build ./xin.go -o ./xin
```

to build Xin as a standalone binary, or run `make install` to install Xin alongside the vim syntax highlighting / definition file.

Xin can currently run as a repl, or execute from a file. To run the repl, simply run

```
xin repl
```

And a prompt will appear. Each input and output line in the prompt is numbered, and you can access previous results from the repl with the corresponding number. For example, take this repl session:

```
0 ) (+ 1 2)
0 ) 3

1 ) (vec 1 2 3)
1 ) (<vec> 1 2 3)
```

We can access `3`, the result from line 0, as the variable `_0`. Likewise, `_1` will reference the vector from line 1.

To run Xin programs from files, use `xin run`. For example, to run `samples/list.xin`

```
xin run samples/list.xin
```

## Key ideas explored

While Xin is meant to be a practical general-purpose programming language, as a toy project, it explores a few key ideas that I couldn't elegantly fit into Ink, my first language.

### Lazy evaluation

In Xin, the evaluation of every expression is deferred until it's absolutely necessary. We call this lazy evaluation. Values are kept in "lazy" states until some outside action or special form coerces a lazy value to resolve to a real value.

### Code as data structure

Like all Lisps, the power and flexibility of Xin comes from its syntax, and the fact that it's trivial to express Xin programs as simple nested list data structures. Xin is not a "real" Lisp, because Xin does not use linked lists internally to represent its own syntax. But the idea of runtime syntax introspection with list data structures is alive in Xin, and it's much of what makes the language so extensible.

Expressing all syntax as lists (auto-resizing vectors, internally) allows functions to modify or introspect their syntax at runtime, and this allows Xin programmers to build language constructs like switch cases, variadic functions, and lambdas in the userspace, as a library function, rather than having to add them to the language core, which remains tiny.

### Streaming I/O and pipes

Xin expresses all I/O operations as operations to streams.

### Syntax minimalism and extensibility

Xin is unusual among even toy programming languages in that there are only four special forms defined in the language spec: `:` (called "bind"), `if`, `do`, and `import`. All other language constructs, like loops, switch cases, and `unless` are defined in the standard library as Xin forms, not in the runtime as special cases.
