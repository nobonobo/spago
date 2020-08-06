package new

import (
	"strings"
	"text/template"
)

const src = `package {{.PkgName}}

import (
	"github.com/nobonobo/spago"
)

//go:generate spago generate -c {{.Name}} -p {{.PkgName}} {{.Name | ToLower}}.html

// {{.Name}}  ...
type {{.Name}} struct {
	spago.Core
}
`

var templ = template.Must(template.New("").Funcs(template.FuncMap{
	"ToUpper": strings.ToUpper,
	"ToLower": strings.ToLower,
}).Parse(src))
