package engine

import (
	"bytes"
	"text/template"
)

type Vars map[string]any

type Engine struct {
	funcs template.FuncMap
}

func New() *Engine {
	return &Engine{
		funcs: funcMap(),
	}
}

func (e *Engine) Render(text string, vars Vars) (string, error) {
	t, err := template.New("").Funcs(e.funcs).Parse(text)
	if err != nil {
		return text, err
	}
	var buf bytes.Buffer

	err = t.Execute(&buf, vars)
	if err != nil {
		return text, err
	}
	return buf.String(), nil
}
