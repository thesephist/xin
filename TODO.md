# Todos

- [ ] Generating generative art postcards with Xin with bmp.xin
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
