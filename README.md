:x: loc
=======

[loc](https://godoc.org/github.com/rhysd/loc) is a small library to make errors work with source
files with location information. It provides a struct to represent a source file, a specific
position in code and an error related to specific range or position in source.

This library is useful to provide a unified look for error messages of compilers, interpreters or
translators.

By using `Source` and `Position` types as position information, this library provides an error
type which provides nice look error message.

- It shows the code snippet causing an error
- Enable to add notes to error by nesting an error instance like [pkg/errors](https://github.com/pkg/errors)
- Proper location is automatically added to error messages and notes
- Colorized label like 'Error:' or 'Note:'

TODO: screenshot

TODO: code snippet
