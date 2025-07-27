# `util`
[![Go Reference](https://pkg.go.dev/badge/go.dsage.org/standard/util.svg)][1]

This package contains utility classes that do not quite fit into other packages but are still useful enough to include.

## `MultiWriter`
MultiWriter is a custom implementation of the `io.MultiWriter` method that creates a single writer which outputs all
written data to all provided `io.Writer` objects. This implementation is different in that you can unregister writers
on the fly.

[1]: https://pkg.go.dev/go.dsage.org/standard/util
