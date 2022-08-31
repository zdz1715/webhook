package engine

import (
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

func funcMap() template.FuncMap {
	f := sprig.TxtFuncMap()

	extra := template.FuncMap{}

	for k, v := range extra {
		f[k] = v
	}
	return f
}
