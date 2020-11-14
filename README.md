# interfacegen

interfacegen is a utility to generate interfaces based on concrete types with methods.

It takes inspiration from  [vburenin/ifacemaker](https://github.com/vburenin/ifacemaker) and [rjeczalik/interfaces](https://github.com/rjeczalik/interfaces).

It leverages [golang.org/x/tools/go/packages](https://godoc.org/golang.org/x/tools/go/packages) to parse packages, and optionally depends on [go/ast](https://golang.org/pkg/go/ast/) for comments.

Its feature set is as follows:

- Supports generating interfaces for all types with methods in a package
- Supports copying comments from the source implementation into interfaces
- Supports unique names for generated interfaces via `-t concreteName,interfaceName`
- Correctly imports types defined in the source package

## Installation

```sh
go get github.com/Deiz/interfacegen
```

## Usage 

### CLI

```sh
# Generate FooInterface from concrete type Foo
interfacegen -s github.com/your/package -t Foo,FooInterface
```

### go generate

```go
// Takes all types with methods in the current package and generates
// interfaces for them.
//go:generate interfacegen -s . -o interfaces/interfaces.go
```
