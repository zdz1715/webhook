package engine

import (
	"net/url"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

func funcMap() template.FuncMap {
	f := sprig.TxtFuncMap()

	extra := template.FuncMap{
		// url
		"urlEncode": urlEncode,
		"urlDecode": urlDecode,
	}

	for k, v := range extra {
		f[k] = v
	}
	return f
}

// urlEncode encodes an item into a url string
func urlEncode(v interface{}) string {
	if s, ok := v.(string); ok {
		return url.QueryEscape(s)
	}
	return ""
}

// urlDecode encodes an item into a url string
func urlDecode(v interface{}) (string, error) {
	if s, ok := v.(string); ok {
		return url.QueryUnescape(s)
	}
	return "", nil
}
