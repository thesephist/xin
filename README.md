# Xin

Xin (신/心) is a simple, fast, evented Lisp-like programming language.

## Goals

- Runtime / async performance
- Expressive, lisp-like functional syntax / semantics
    - Proper tail calls
- Portability to resource-constrained environments (within constraints of Go)
- Hygienic macros
- Bytecode compilation and export
- Embeddability in other Go applications
- Programs that lend themselves to sophisticated data structures and dumb algorithms, rather than vice versa
- Great REPL experience
- Ease of learning, within the constraints of the syntax

Non-goals

- Suitability in microcontrollers / non-desktop/server architectures
- Competing in performance against C / fully compiled languages

## Implementation wishes

- Xin is lazily evaluated. (or at least, tries to be, but this is best effort.)

- Xin's sole branching primitive is the pattern-matching `::`, borrowed from Ink. Other conditionals like `if`, `unless` are composed from the match function.
- Xin's sole looping construct is recursion. Other constructs like the `for` and `range` functions are composed from recursion, potentially in macros.
- Xin should have a native TCP/ UDP `net` interface, rather than a native HTTP adapter like Ink neglecting raw network syscalls.

- The way we achieve this is by taking a prefix syntax and translating it into 100% postfix syntax before compiling it into machine-portable bytecode, which runs in a stack machine that can execute postfix functions quickly.
- To make Xin lisp-style powerful, we need some way to define macros or give quoted expressions to functions that are not evaluated immediately.
- Ink is about language expressiveness and maximum flexibility that reflected real-world programs. Xin is about a virtual machine and data structures forming a data pipeline.
- Xin will be written in Go, but maybe C and Rust later?
- For performance's sake Xin tries to minimize Gc in the interpreter and language to zero. This means Xin is not memory managed, though we provide APIs to allocate and get a reference to a `Arc` of a memory. Not a pointer, but a fixed-size reference. And we can then also practice writing a memory allocator maybe.
- Xin will be single-threaded and evented / non-blocking, using callbacks.
- Xin will have a more sophisticated REPL than Ink, because it's allowed to have dependencies on third party code as it's not mission critical.
- Another focus of Xin is compiler optimization, and the language should allow for ample compiler optimizations.
- Xin should keep track of stack and return rich, valuable runtime error traces.
- Xin supports proper tail calls (is tail call optimized) at zero runtime cost, unlike Ink, which has a runtime cost of a Go function call, because a bytecode instruction accounts for a proper tail call.

## Syntax brainstorming

Lisp-style

- Also reference: https://schemers.org/Documents/Standards/R5RS/r5rs.pdf

```
-- Almost every symbol is defined in userland
-- except for those directly linked to machine instructions
-- and native function calls.
(: (f x y z) (+ x (+ y z)))
(: (sq n) (* n n))
```

## Type system

Supports algebraic types and type aliases, but not parameterized types.

Maybe if we have a macro system we can use that a stand-in for generics.

Primitive types
- `String`, which also encodes any binary blob data
- `Int`, word-size integer
- `Float`, 64-bit float
- `Bool`, boolean
- `List`, lisp-style dynamic array

## Xin bytecode

Xin bytecode is a stack machine that is the compilation target for the Xin language.

- Look at WebAssembly bytecode spec.
