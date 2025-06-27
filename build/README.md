# `build`
[![Go Reference](https://pkg.go.dev/badge/go.dsage.org/standard/build.svg)][1]

This package contains logic to identify the exact revision of the app that is running by reading the VCS information
stored within the binary at compile time. The VCS information can be ignored in favor of a statically set version if one
is embedded manually. Whether using VCS information or a static version, the version will be stored in `DetectedVersion`
during initialization. The detected Go version will also be stored in `GoVersion` if you prefer to access it using this
package as well.

## VCS Build Info
Upon initialization, the package will automatically read the embedded build info to identify the VCS revision and state
of the repository when the app was built. The detected app version is stored in `DetectedVersion` and the detected Go
version is stored in `GoVersion`. If no app version can be identified, it will be set to `dev`. If the repository was
not clean (any non-ignored files were changed) the version will end with `-dirty`.

### Example
```go
package main

import (
  "fmt"

  "go.dsage.org/standard/build"
)

func main() {
	// compiled with Go 1.24.3
	fmt.Println("The Go version is", build.GoVersion)

	// with no VCS information stored
	fmt.Println("The app version is", build.DetectedVersion)

	// with VCS commit ID `1234567890`
	fmt.Println("The app version is", build.DetectedVersion)

	// with VCS commit ID `1234567890` and modified files
	fmt.Println("The app version is", build.DetectedVersion)

	// Output:
	// The Go version is go1.24.3
	// The app version is dev
	// The app version is 1234567890
	// The app version is 1234567890-dirty
}
```

## Static
If the app is being built as part of a release process with a specific version that should be displayed, that version
can be embedded at compile time to be displayed instead of any VCS information. This can be done automatically with any
build system that allows you to add `ldflags` to the `go build` command.

### Example
```shell
go build -ldflags "-X go.dsage.org/standard/build.staticVersion=1.2.34" <app directory>
```
```go
package main

import (
  "fmt"

  "go.dsage.org/standard/build"
)

func main() {
	// compiled with Go 1.24.3
	fmt.Println("The Go version is", build.GoVersion)

	// with a `staticVersion` of `1.2.34` embedded at compile time
	fmt.Println("The app version is", build.DetectedVersion)

	// Output:
	// The Go version is go1.24.3
	// The app version is 1.2.34
}
```

[1]: https://pkg.go.dev/go.dsage.org/standard/build
