package plugin

import (
	"fmt"

	"github.com/Neakxs/protoc-gen-authz/authorize"
	"github.com/Neakxs/protoc-gen-authz/internal/cfg"
	"github.com/Neakxs/protoc-gen-authz/internal/template"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/interpreter"
	v1alpha1 "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

func NewFile(p *protogen.Plugin, f *protogen.File, c *cfg.Config) *File {
	g := p.NewGeneratedFile(f.GeneratedFilenamePrefix+".pb.authz.go", f.GoImportPath)
	msgs := []*Message{}
	for i := 0; i < len(f.Messages); i++ {
		msgs = append(msgs, NewMessage(f.Messages[i], c))
	}
	return &File{
		p:        p,
		g:        g,
		File:     f,
		Config:   c,
		Messages: msgs,
	}
}

type File struct {
	p *protogen.Plugin
	g *protogen.GeneratedFile
	*protogen.File
	Config   *cfg.Config
	Messages []*Message
}

func (f *File) Generate() error {
	if err := f.Validate(); err != nil {
		return err
	}
	if tmpl, err := template.GenerateTemplate(f.p.Request.CompilerVersion, f.g); err != nil {
		return err
	} else {
		return tmpl.Execute(f.g, f)
	}
}

func (f *File) Validate() error {
	for i := 0; i < len(f.Messages); i++ {
		if err := f.Messages[i].Validate(); err != nil {
			return err
		}
	}
	return nil
}

func NewMessage(m *protogen.Message, c *cfg.Config) *Message {
	return &Message{Message: m, Config: c}
}

type Message struct {
	*protogen.Message
	Config *cfg.Config
}

func (m *Message) MessageRule() *authorize.Rule {
	return proto.GetExtension(m.Desc.Options(), authorize.E_Rule).(*authorize.Rule)
}

func (m *Message) Macros() ([]string, error) {
	if m.MessageRule() == nil {
		return nil, nil
	}
	envOpts := []cel.EnvOption{
		cel.Types(&authorize.AuthorizationContext{}),
		cel.DeclareContextProto(m.Desc),
		cel.Declarations(
			decls.NewVar("_ctx", decls.NewObjectType(string((&authorize.AuthorizationContext{}).ProtoReflect().Descriptor().FullName()))),
		),
	}
	for k := range m.Config.Globals.Functions {
		envOpts = append(envOpts, cel.Declarations(decls.NewFunction(k, decls.NewOverload(k, []*v1alpha1.Type{}, &v1alpha1.Type{TypeKind: &v1alpha1.Type_Primitive{Primitive: v1alpha1.Type_BOOL}}))))
	}
	env, err := cel.NewEnv(envOpts...)
	if err != nil {
		return nil, err
	}
	ast, issues := env.Compile(m.MessageRule().Expr)
	if issues != nil && issues.Err() != nil {
		return nil, issues.Err()
	}
	return findMacrosInAST(ast, m.Config.Globals.Functions), nil
}

func (m *Message) Validate() error {
	if m.MessageRule() == nil {
		return nil
	}
	mm := map[string]string{}
	if macros, err := m.Macros(); err != nil {
		return err
	} else {
		for _, macro := range macros {
			mm[macro] = m.Config.Globals.Functions[macro]
		}
	}
	envOpts, err := authorize.BuildEnvOptionsWithMacros([]cel.EnvOption{
		cel.Types(&authorize.AuthorizationContext{}),
		cel.DeclareContextProto(m.Desc),
		cel.Declarations(
			decls.NewVar("_ctx", decls.NewObjectType(string((&authorize.AuthorizationContext{}).ProtoReflect().Descriptor().FullName()))),
		),
	}, mm)
	if err != nil {
		return err
	}
	env, err := cel.NewEnv(envOpts...)
	if err != nil {
		return err
	}
	ast, issues := env.Compile(m.MessageRule().Expr)
	if issues != nil && issues.Err() != nil {
		return issues.Err()
	}
	switch ast.ResultType().TypeKind.(type) {
	case *v1alpha1.Type_Primitive:
		if v1alpha1.Type_PrimitiveType(ast.ResultType().GetPrimitive().Number()) != v1alpha1.Type_BOOL {
			return fmt.Errorf("result type not bool")
		}
	default:
		return fmt.Errorf("result type not bool")
	}
	_, err = env.Program(ast, cel.OptimizeRegex(interpreter.MatchesRegexOptimization))
	if err != nil {
		return err
	}
	return nil
}
