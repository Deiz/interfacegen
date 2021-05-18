package main

import (
	"bytes"
	"text/template"
)

type methodDef struct {
	Method string
	Slug   string
	Doc    []string
}

type interfaceDef struct {
	Name    string
	Doc     []string
	Methods []methodDef
}

const packageTemplate = `
{{ if .Comment -}}
// {{ .Comment -}}
{{- end }}

package {{.PackageName}}

{{ if .Imports -}}
import (
{{- range $_, $import := .Imports }}
	{{ $import | printf "%q" }}
{{- end }}
)
{{ end -}}

{{ range $_, $interface := .Interfaces }}
{{- range $_, $doc := $interface.Doc }}
{{ $doc }}
{{- end }}
type {{ $interface.Name }} interface {
{{- range $_, $method := $interface.Methods }}
{{- range $_, $doc := $method.Doc }}
	{{ $doc }}
{{- end }}
	{{ $method.Method }}
{{- end }}
}
{{ end }}
`

var (
	tmpl = template.Must(template.New("interface").Parse(packageTemplate))
)

func (app *application) generatePackage(imports []string, interfaceDefs []interfaceDef) (code string, err error) {
	buf := bytes.NewBuffer(nil)

	err = tmpl.Execute(buf, struct {
		Comment     string
		PackageName string
		Imports     []string
		Interfaces  []interfaceDef
	}{
		Comment:     app.Comment,
		PackageName: app.DstPackage,
		Imports:     imports,
		Interfaces:  interfaceDefs,
	})
	if err != nil {
		return code, err
	}

	return buf.String(), err
}
