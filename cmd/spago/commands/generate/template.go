package generate

import "text/template"

var templ = template.Must(template.New("").Parse(`package {{.PkgName}}

import (
{{range $v, $_ := .StdImports}}{{printf "\t%q\n" $v}}{{end -}}
{{if .StdImports}}{{printf "\n"}}{{end -}}
{{- range $v, $as := .Imports -}}
  {{- $length := len $as -}}
	{{- if gt $length 0 -}}
		{{- printf "\t$s %q\n" $as $v -}}
	{{- else -}}
		{{- printf "\t%q\n" $v -}}
	{{- end -}}
{{- end -}}
)

// Render ...
func (c *{{.ComponentName}}) Render() spago.HTML {
	return {{- .Generated}}
}

{{- range $_, $code := .AppendCode}}
{{$code}}
{{- end}}
`))
