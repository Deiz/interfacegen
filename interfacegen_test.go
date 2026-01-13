package main

import (
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

	test5Input = `package foo

import (
	r8 "github.com/go-redis/redis/v8"
	r9 "github.com/redis/go-redis/v9"

	"golang.org/x/tools/imports"

	_ "embed"
)

func main() {
}

//go:embed response_codes.json
var responseCodes []byte

type FooClient struct {
	v8 *r8.Client
	v9 *r9.Client
}

type Str string

func (f *FooClient) ThingA() (Str, *r8.BoolCmd) {
	return f.v8.DoThingA()
}

func (f *FooClient) ThingB() (*r9.BoolCmd, imports.Options) {
	return f.v9.DoThingB()
}
`

	test5Output = `package interfaces

import (
	"foo"

	r8 "github.com/go-redis/redis/v8"
	r9 "github.com/redis/go-redis/v9"
	"golang.org/x/tools/imports"
)

type FooClient interface {
	ThingA() (foo.Str, *r8.BoolCmd)
	ThingB() (*r9.BoolCmd, imports.Options)
}
`

	goMod = `module foo

go 1.16`

	goModRedis = `module foo

go 1.23.0

require (
	github.com/go-redis/redis/v8 v8.11.5
	github.com/redis/go-redis/v9 v9.12.0
	golang.org/x/tools v0.36.0
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	golang.org/x/mod v0.27.0 // indirect
	golang.org/x/sync v0.16.0 // indirect
)
`

	goSumRedis = `
github.com/bsm/ginkgo/v2 v2.12.0 h1:Ny8MWAHyOepLGlLKYmXG4IEkioBysk6GpaRTLC8zwWs=
github.com/bsm/ginkgo/v2 v2.12.0/go.mod h1:SwYbGRRDovPVboqFv0tPTcG1sN61LM1Z4ARdbAV9g4c=
github.com/bsm/gomega v1.27.10 h1:yeMWxP2pV2fG3FgAODIY8EiRE3dy0aeFYt4l7wh6yKA=
github.com/bsm/gomega v1.27.10/go.mod h1:JyEr/xRbxbtgWNi8tIEVPUYZ5Dzef52k01W3YH0H+O0=
github.com/cespare/xxhash/v2 v2.3.0 h1:UL815xU9SqsFlibzuggzjXhog7bL6oX9BbNZnL2UFvs=
github.com/cespare/xxhash/v2 v2.3.0/go.mod h1:VGX0DQ3Q6kWi7AoAeZDth3/j3BFtOZR5XLFGgcrjCOs=
github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f h1:lO4WD4F/rVNCu3HqELle0jiPLLBs70cWOduZpkS1E78=
github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f/go.mod h1:cuUVRXasLTGF7a8hSLbxyZXjz+1KgoB3wDUb6vlszIc=
github.com/fsnotify/fsnotify v1.4.9 h1:hsms1Qyu0jgnwNXIxa+/V/PDsU6CfLf6CNO8H7IWoS4=
github.com/fsnotify/fsnotify v1.4.9/go.mod h1:znqG4EE+3YCdAaPaxE2ZRY/06pZUdp0tY4IgpuI1SZQ=
github.com/go-redis/redis/v8 v8.11.5 h1:AcZZR7igkdvfVmQTPnu9WE37LRrO/YrBH5zWyjDC0oI=
github.com/go-redis/redis/v8 v8.11.5/go.mod h1:gREzHqY1hg6oD9ngVRbLStwAWKhA0FEgq8Jd4h5lpwo=
github.com/google/go-cmp v0.6.0 h1:ofyhxvXcZhMsU5ulbFiLKl/XBFqE1GSq7atu8tAmTRI=
github.com/google/go-cmp v0.6.0/go.mod h1:17dUlkBOakJ0+DkrSSNjCkIjxS6bF9zb3elmeNGIjoY=
github.com/nxadm/tail v1.4.8 h1:nPr65rt6Y5JFSKQO7qToXr7pePgD6Gwiw05lkbyAQTE=
github.com/nxadm/tail v1.4.8/go.mod h1:+ncqLTQzXmGhMZNUePPaPqPvBxHAIsmXswZKocGu+AU=
github.com/onsi/ginkgo v1.16.5 h1:8xi0RTUf59SOSfEtZMvwTvXYMzG4gV23XVHOZiXNtnE=
github.com/onsi/ginkgo v1.16.5/go.mod h1:+E8gABHa3K6zRBolWtd+ROzc/U5bkGt0FwiG042wbpU=
github.com/onsi/gomega v1.18.1 h1:M1GfJqGRrBrrGGsbxzV5dqM2U2ApXefZCQpkukxYRLE=
github.com/onsi/gomega v1.18.1/go.mod h1:0q+aL8jAiMXy9hbwj2mr5GziHiwhAIQpFmmtT5hitRs=
github.com/redis/go-redis/v9 v9.12.0 h1:XlVPGlflh4nxfhsNXPA8Qp6EmEfTo0rp8oaBzPipXnU=
github.com/redis/go-redis/v9 v9.12.0/go.mod h1:huWgSWd8mW6+m0VPhJjSSQ+d6Nh1VICQ6Q5lHuCH/Iw=
golang.org/x/mod v0.27.0 h1:kb+q2PyFnEADO2IEF935ehFUXlWiNjJWtRNgBLSfbxQ=
golang.org/x/mod v0.27.0/go.mod h1:rWI627Fq0DEoudcK+MBkNkCe0EetEaDSwJJkCcjpazc=
golang.org/x/net v0.43.0 h1:lat02VYK2j4aLzMzecihNvTlJNQUq316m2Mr9rnM6YE=
golang.org/x/net v0.43.0/go.mod h1:vhO1fvI4dGsIjh73sWfUVjj3N7CA9WkKJNQm2svM6Jg=
golang.org/x/sync v0.16.0 h1:ycBJEhp9p4vXvUZNszeOq0kGTPghopOL8q0fq3vstxw=
golang.org/x/sync v0.16.0/go.mod h1:1dzgHSNfp02xaA81J2MS99Qcpr2w7fw1gpm99rleRqA=
golang.org/x/sys v0.35.0 h1:vz1N37gP5bs89s7He8XuIYXpyY0+QlsKmzipCbUtyxI=
golang.org/x/sys v0.35.0/go.mod h1:BJP2sWEmIv4KK5OTEluFJCKSidICx8ciO85XgH3Ak8k=
golang.org/x/text v0.3.6 h1:aRYxNxv6iGQlyVaZmk6ZgYEDa+Jg18DxebPSrd6bg1M=
golang.org/x/text v0.3.6/go.mod h1:5Zoc/QRtKVWzQhOtBMvqHzDpF6irO9z98xDceosuGiQ=
golang.org/x/tools v0.36.0 h1:kWS0uv/zsvHEle1LbV5LE8QujrxB3wfQyxHfhOk0Qkg=
golang.org/x/tools v0.36.0/go.mod h1:WBDiHKJK8YgLHlcQPYQzNCkUxUypCaa5ZegCVutKm+s=
gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 h1:uRGJdciOHaEIrze2W8Q3AKkepLTh2hOroT7a+7czfdQ=
gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7/go.mod h1:dt/ZhP58zS4L8KSrWDmTeBkI65Dw0HsyUHuEVlX15mw=
gopkg.in/yaml.v2 v2.4.0 h1:D8xgwECY7CYvx+Y2n4sBz93Jn9JRvxdiyyo8CTfuKaY=
gopkg.in/yaml.v2 v2.4.0/go.mod h1:RDklbk79AGWmwhnvt/jBztapEOGDOx6ZbXqjP6csGnQ=
`
)

func TestInterfacegen(t *testing.T) {
	tests := map[string]struct {
		input  string
		goMod  string
		goSum  string
		output string
		app    *application
	}{
		"happy path": {
			input:  test1Input,
			goMod:  goMod,
			output: test1Output,
		},
		"type skip": {
			input:  test2Input,
			goMod:  goMod,
			output: test2Output,
		},
		"no comments": {
			input:  test3Input,
			goMod:  goMod,
			output: test3Output,
			app: &application{
				IncludeDocs: false,
			},
		},
		"method skip": {
			input:  test4Input,
			goMod:  goMod,
			output: test4Output,
		},
		"conflicting imports": {
			input:  test5Input,
			goMod:  goModRedis,
			goSum:  goSumRedis,
			output: test5Output,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			require := require.New(t)

			inputDir, err := os.MkdirTemp("", "")
			require.NoError(err)

			defer func() {
				//os.RemoveAll(inputDir)
			}()

			t.Logf("Created %s to store temporary input", inputDir)

			require.NoError(os.Chdir(inputDir))
			require.NoError(os.WriteFile(filepath.Join(inputDir, "main.go"), []byte(test.input), 0644))
			require.NoError(os.WriteFile(filepath.Join(inputDir, "go.mod"), []byte(test.goMod), 0644))
			require.NoError(os.WriteFile(filepath.Join(inputDir, "go.sum"), []byte(test.goSum), 0644))

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

			require.NoError(app.Run())

			data, err := os.ReadFile(outputFile.Name())
			require.NoError(err)

			require.Equal(test.output, string(data))
			require.NoError(os.Remove(outputFile.Name()))
		})
	}
}
