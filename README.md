# Standard
[![Go Reference](https://pkg.go.dev/badge/go.dsage.org/standard.svg)][1]

Standard is a small set of packages that are shared between many of my projects, both internal and public. These have
been broken out into their own project to eliminate the need to reimplement them or apply fixes to many projects at once
if an issue is discovered.

While the repository is public, the primary focus is my own use and thus responses may be slow.

## Getting Started
Start by adding the package as a dependency to your project:

```shell
go get -u go.dsage.org/standard
```

For an overview of each package within the library, see the [documentation][1].

## Requirements
There are no direct dependencies for this project, but testing only takes place on the following hosts:

- Linux
  - Arch Linux (manual)
  - Ubuntu (automated, via GitHub Actions: `ubuntu-latest`)
- macOS (automated, via GitHub Actions: `macos-latest`)
- Windows (automated, via GitHub Actions: `windows-latest`)

## Security Vulnerabilities
Please see the separate security policy in [`SECURITY.md`](SECURITY.md).

## License
Standard is released under the [MIT License](https://choosealicense.com/licenses/mit/) (see [`LICENSE`](LICENSE)).

[1]: https://pkg.go.dev/go.dsage.org/standard
