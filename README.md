# Xin

Xin is a purely functional programming language inspired by lisp syntax and CSP ideas about concurrency and data sharing. Xin aspires to be an expressive, concise language built on a small number of simple elements that work well together.

## Goals

- Expressive, extensible syntax that's natural to read and suitable for defining DSLs
- Performance: [Ink](https://github.com/thesephist/ink) is about 3-4x slower than Python. Xin programs should run at least 10x faster than Python, and the language design should lend itself to efficient compiler optimizations.
- Ability to define special forms on the fly through lazy evaluation of forms, but not macros.
    - i.e. `(if cond ifTrue if False)` is not a builtin, and should be defined by the standard library based on `(match v ...)`.
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
