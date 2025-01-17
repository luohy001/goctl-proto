package proto

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
	"text/template"
	"unicode"
)

var fileTemplate = `syntax = "{{ .Syntax }}";

package {{ replace .Package "-" "." }};

{{- range .Options }}
option {{ .Name }} = "{{ .Value }}";
{{- end }}

{{- range .Messages }}
{{range .Descs }}
// {{ . }}
{{- end }}
{{- $fieldCount := len .Fields }}
message {{ title .Name }} { {{- if eq $fieldCount 0 -}} }{{ else }}
{{- range $index,$value := .Fields }}
{{- range .Descs }}
    // {{ . }}
{{- end }}
{{- $tagCount := len .Tags }}
    {{ if .Repeated }}repeated {{ end }}{{ .TypeName }} {{ .Name }} = {{ add $index 1 }}{{- if eq $tagCount 0 -}};{{ else }} [{{ join .Tags ", " }}];{{- end }}
{{- end }}
}
{{- end }}
{{- end }}

{{- range .Services }}
{{ range .Descs }}
// {{ . }}
{{- end }}
service {{ title .Name }} {
{{- range .Rpcs }}
{{- range .Descs }}
    // {{ . }}
{{- end }}
    rpc {{ title .Name }} ({{ if ne .Request nil }}{{ .Request.Name }}{{ else }}Empty{{ end }}) returns ({{ if ne .Response nil }}{{ .Response.Name }}{{ else }}Empty{{ end }});
{{- end }}
}
{{- end }}
`

var funcMap = template.FuncMap{
	"add": func(x, y int) int {
		return x + y
	},
	"title": func() func(string) string {
		caser := cases.Title(language.English, cases.NoLower)
		return func(s string) string {
			ss := strings.FieldsFunc(s, func(c rune) bool {
				return !unicode.IsLetter(c) && !unicode.IsNumber(c)
			})
			for i := range ss {
				ss[i] = caser.String(ss[i])
			}
			return strings.Join(ss, "")
		}
	}(),
	"replace": func(s, old, new string) string {
		return strings.ReplaceAll(s, old, new)
	},
	"join": func(s []string, sep string) string {
		return strings.Join(s, sep)
	},
}
