# Todos

- [ ] Rich runtime error traces and suggestions for corrections like Rust.
- [ ] Well thought-out macro system. Define with `::` pass AST as `vec`s.
    - `MacroFormValue`
    - We need two utility functions, one to convert from `[]*astNode` to `VecValue`, and another to go the other way.
    - The main problem here is that we don't have a good way of propagating position and stack trace data thru the macro system.
    - I think we need a native AST representation datatypes to do this right, rather than trying to shoehorn AST data into primitive types.
- [ ] Add date / time native forms. Unix to ISO 8601 should be a built-in for now. Might move to stdlib later.
- [ ] Native TCP / UCP `os::net` interfaces, on top of which an HTTP library can be written.
- [ ] More robust filesystem IO APIs in general -- supporting read/write offsets, permissions
    - Consider making filesystem I/O asynchronous. In particular, stream sink/source on files, and `os::delete`
    - Update os::open to allow file to be truncated / appended at an offset and to read from offsets.
    - Reading reads 4K bytes though the buffer size can also be modified as a side effect.
