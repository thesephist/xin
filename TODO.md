# Todos

- [ ] BMP library -- for creating lisp art.
    - Generating generative art postcards with Xin
- [ ] Runtime stack traces for errors
- [ ] Well thought-out macro system. Define with `::` pass AST as `vec`s.
    - `MacroFormValue`
    - We need two utility functions, one to convert from `[]*astNode` to `VecValue`, and another to go the other way.
    - The main problem here is that we don't have a good way of propagating position and stack trace data thru the macro system.
    - I think we need a native AST representation datatypes to do this right, rather than trying to shoehorn AST data into primitive types.
- [ ] Xin should incorporate an intermediate representation that reflects and allows for lots of static analysis. Tending towards a compiler.
    - Static analysis:
        - Statically resolve references, since Xin is lexically scoped all of the time
        - Statically determine places where lazy evaluation has no benefit, and don't lazy-evaluate (remove indirection)
        - Inline small functions ("small" here is probably the number of nodes in the AST)
        - Statically determine object / value lifetimes and maybe deterministically allocate memory for those slots
- [ ] Native TCP / UCP `net` interfaces, on top of which an HTTP library can be written.
    - `(os::open <socket>)` like `(os::open '127.0.0.1' 8080)` will open a TCP/UDP connection stream. The stream will emit streams which correspond to individual connections (reader & writer streams).
