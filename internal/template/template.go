package template

import (
	_ "embed"
	"fmt"
	"text/template"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

//go:embed template.go.tmpl
var tmpl string

func GenerateTemplate(v *pluginpb.Version, g *protogen.GeneratedFile) (*template.Template, error) {
	return template.New("").Funcs(template.FuncMap{
		"PluginVersion": func() string { return "v0.0.0" },
		"ProtocVersion": func() string {
			return fmt.Sprintf("v%d.%d.%d", *v.Major, *v.Minor, *v.Patch)
		},
	}).Funcs(template.FuncMap{
		"context": func(s string) string {
			return g.QualifiedGoIdent(protogen.GoImportPath("context").Ident(s))
		},
		"sync": func(s string) string {
			return g.QualifiedGoIdent(protogen.GoImportPath("sync").Ident(s))
		},
		"authorize": func(s string) string {
			return g.QualifiedGoIdent(protogen.GoImportPath("github.com/Neakxs/protoc-gen-authz/authorize").Ident(s))
		},
		"cel": func(s string) string {
			return g.QualifiedGoIdent(protogen.GoImportPath("github.com/google/cel-go/cel").Ident(s))
		},
		"decls": func(s string) string {
			return g.QualifiedGoIdent(protogen.GoImportPath("github.com/google/cel-go/checker/decls").Ident(s))
		},
		"interpreter": func(s string) string {
			return g.QualifiedGoIdent(protogen.GoImportPath("github.com/google/cel-go/interpreter").Ident(s))
		},
		"parser": func(s string) string {
			return g.QualifiedGoIdent(protogen.GoImportPath("github.com/google/cel-go/parser").Ident(s))
		},
		"codes": func(s string) string {
			return g.QualifiedGoIdent(protogen.GoImportPath("google.golang.org/grpc/codes").Ident(s))
		},
		"status": func(s string) string {
			return g.QualifiedGoIdent(protogen.GoImportPath("google.golang.org/grpc/status").Ident(s))
		},
	}).Parse(tmpl)
}
