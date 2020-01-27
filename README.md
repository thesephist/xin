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
- Great REPL

## Key ideas in Xin

[section to be written]
