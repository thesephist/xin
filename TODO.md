# Todos

- [ ] Native TCP / UCP `os::net` interfaces, on top of which an HTTP library can be written.
- [ ] Rich runtime error traces and suggestions for corrections like Rust.
- [ ] More robust filesystem IO APIs in general -- supporting read/write offsets, permissions
    - [ ] Consider making filesystem I/O asynchronous. In particular, stream sink/source on files, and `os::delete`
