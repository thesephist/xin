# Xin programming language

[![GoDoc](https://godoc.org/github.com/thesephist/xin?status.svg)](https://godoc.org/github.com/thesephist/xin)
[![Build Status](https://travis-ci.com/thesephist/xin.svg?branch=master)](https://travis-ci.com/thesephist/xin)


Xin is a purely functional programming language inspired by lisp syntax and CSP ideas about concurrency and data sharing. Xin aspires to be an expressive, concise language built on a small number of simple elements that work well together.

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
- Performance: [Ink](https://github.com/thesephist/ink) is about 3-4x slower than Python. Xin programs should run at least 10x faster than Python, and the language design should lend itself to efficient compiler optimizations.
- Programs that lend themselves to sophisticated data structures and dumb algorithms, rather than vice versa
- Ease of learning, within the constraints of the syntax
- Great REPL

## Syntax (WIP)

- Also reference: https://schemers.org/Documents/Standards/R5RS/r5rs.pdf

## Type system (WIP)

Supports algebraic types and type aliases and parameterized types, with strong type inference.

Maybe if we have a macro system we can use that a stand-in for generics.

## Xin bytecode

Xin bytecode is a stack machine that is the compilation target for the Xin language.

- Look at WebAssembly bytecode spec.
- Look at ALE (https://github.com/kode4food/ale) implementation and syntax, especially for the stack machine bytecode vm
