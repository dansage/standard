# `env`
This package contains a very simple environment detection system. In most situations, this package can be used simply by
reading the contents of the `Current` variable to determine the current environment type. If the environment type is not
explicitly set or the value of the `ENV` environment variable is unknown, `Production` is assumed.

## Automatic Detection
Upon initialization the package will automatically use the contents of the `ENV` environment variable to identify the
current environment type and store it in the `Current` variable.

### Example
```go
package main

import (
  "fmt"

  "go.dsage.org/standard/env"
)

func main() {
	// with `ENV` set to `test` or `testing
	fmt.Println("The current environment is", env.Current)

	// with `ENV` set to `dev`, `development`, or `local`
	fmt.Println("The current environment is", env.Current)

	// with `ENV` unset or set to anything but the above values
	fmt.Println("The current environment is", env.Current)

	// Output:
	// The current environment is testing
	// The current environment is development
	// The current environment is production
}
```

## Manual Detection
If the environment variables are changed after the package is initialized, you can manually detect the current
environment type with `DetectEnvironment()`. This method does _not_ update `Current` automatically.

### Example
```go
package main

import (
  "fmt"
  "os"

  "go.dsage.org/standard/env"
)

func main() {
	// with `ENV` set to `test` or `testing
	fmt.Println("The current environment is", env.Current)

	// change the contents of the `ENV` variable to indicate this is a development environment
	_ = os.Setenv("ENV", "development")

	// detect the current environment type manually
	current := env.DetectEnvironment()
	fmt.Println("The current environment is", current)

	// Output:
	// The current environment is testing
	// The current environment is development
}
```
