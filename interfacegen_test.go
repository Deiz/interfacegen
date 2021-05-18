package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	test1Input = `package main

func main() {
}

// Hello, world!
type Foo struct {}

// Name prints Foo's name
func (f *Foo) Name() string {
	return "foo"
}
`
	test1Output = `package interfaces

// Hello, world!
type Foo interface {
	// Name prints Foo's name
	Name() string
}
`
	test2Input = `package interfaces

//interfacegen:skip
type Foo struct {}

func (f *Foo) Name() string {
	return "foo"
}

type Bar struct {}

func (b *Bar) Age() int {
	return 3
}
`
	test2Output = `package interfaces

type Bar interface {
	Age() int
}
`
	test3Input = `package main

func main() {
}

// Hello, world!
type Foo struct {}

// Name prints Foo's name
func (f *Foo) Name() string {
	return "foo"
}
`
	test3Output = `package interfaces

type Foo interface {
	Name() string
}
`
	test4Input = `package interfaces

type Foo int

func (f Foo) Val() int {
	return f
}

// interfacegen:skip
func (f *Foo) Incr() {
	*f++
}
`
	test4Output = `package interfaces

type Foo interface {
	Val() int
}
`
	goMod = `module foo

go 1.16`
)

func TestInterfacegen(t *testing.T) {
	tests := map[string]struct {
		input  string
		output string
		app    *application
	}{
		"happy path": {
			input:  test1Input,
			output: test1Output,
		},
		"type skip": {
			input:  test2Input,
			output: test2Output,
		},
		"no comments": {
			input:  test3Input,
			output: test3Output,
			app: &application{
				IncludeDocs: false,
			},
		},
		"method skip": {
			input:  test4Input,
			output: test4Output,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			require := require.New(t)

			inputDir, err := os.MkdirTemp("", "")
			require.NoError(err)

			defer func() {
				os.RemoveAll(inputDir)
			}()

			t.Logf("Created %s to store temporary input", inputDir)

			require.NoError(os.Chdir(inputDir))
			require.NoError(os.WriteFile(filepath.Join(inputDir, "main.go"), []byte(test.input), 0644))
			require.NoError(os.WriteFile(filepath.Join(inputDir, "go.mod"), []byte(goMod), 0644))

			outputFile, err := os.CreateTemp("", "")
			require.NoError(err)
			require.NoError(outputFile.Close())

			t.Logf("Created %s to store temporary output", outputFile.Name())

			var app application

			if test.app == nil {
				app = application{
					IncludeDocs:        true,
					IncludeAllPackages: true,
				}
			} else {
				app = *test.app
			}

			app.SrcPackage = "./"
			app.DstPackage = "interfaces"
			app.Output = outputFile.Name()

			require.NoError(app.Run(context.Background()))

			data, err := os.ReadFile(outputFile.Name())
			require.NoError(err)

			require.Equal(test.output, string(data))
			require.NoError(os.Remove(outputFile.Name()))
		})
	}
}
