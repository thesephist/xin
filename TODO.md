# Todos

- [ ] Add math Rand, crypto Rand, and date / time default forms. Unix to ISO 8601 should be a built-in for now. Might move to stdlib later.
- [ ] Native TCP / UCP `os::net` interfaces, on top of which an HTTP library can be written.
- [ ] Rich runtime error traces and suggestions for corrections like Rust.
- [ ] More robust filesystem IO APIs in general -- supporting read/write offsets, permissions
    - Consider making filesystem I/O asynchronous. In particular, stream sink/source on files, and `os::delete`
    - Update os::open to allow file to be truncated / appended at an offset and to read from offsets.
    - Reading reads 4K bytes though the buffer size can also be modified as a side effect.
