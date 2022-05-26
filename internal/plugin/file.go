package plugin

import (
	"github.com/Neakxs/protoc-gen-authz/authorize"
	"github.com/Neakxs/protoc-gen-authz/cfg"
	"github.com/Neakxs/protoc-gen-authz/internal/template"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

func NewFile(p *protogen.Plugin, f *protogen.File, c *cfg.Config) *File {
	g := p.NewGeneratedFile(f.GeneratedFilenamePrefix+".pb.authz.go", f.GoImportPath)
	// msgs := []*Message{}
	// for i := 0; i < len(f.Messages); i++ {
	// 	msgs = append(msgs, NewMessage(f.Messages[i], c))
	// }
	svcs := []*Service{}
	for i := 0; i < len(f.Services); i++ {
		svcs = append(svcs, NewService(f.Services[i], c))
	}
	return &File{
		p:      p,
		g:      g,
		File:   f,
		Config: c,
		// Messages: msgs,
		Services: svcs,
	}
}

type File struct {
	p *protogen.Plugin
	g *protogen.GeneratedFile
	*protogen.File
	Config   *cfg.Config
	Services []*Service
	// Messages []*Message
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
	for i := 0; i < len(f.Services); i++ {
		if err := f.Services[i].Validate(); err != nil {
			return err
		}
	}
	return nil
}

func NewService(s *protogen.Service, c *cfg.Config) *Service {
	mths := []*Method{}
	for i := 0; i < len(s.Methods); i++ {
		mths = append(mths, NewMethod(s.Methods[i], c))
	}
	return &Service{Service: s, Config: c, Methods: mths}
}

type Service struct {
	*protogen.Service
	Config  *cfg.Config
	Methods []*Method
}

func (s *Service) Validate() error {
	for i := 0; i < len(s.Methods); i++ {
		if err := s.Methods[i].Validate(); err != nil {
			return err
		}
	}
	return nil
}

func NewMethod(m *protogen.Method, c *cfg.Config) *Method {
	return &Method{Method: m, Config: c}
}

type Method struct {
	*protogen.Method
	Config *cfg.Config
}

func (m *Method) MethodRule() *authorize.MethodRule {
	return proto.GetExtension(m.Desc.Options(), authorize.E_Method).(*authorize.MethodRule)
}

func (m *Method) Validate() error {
	if m.MethodRule() == nil {
		return nil
	}
	if _, err := authorize.BuildAuthzProgram(m.MethodRule().GetExpr(), m.Input.Desc, m.Config); err != nil {
		return err
	}
	return nil
}

// func NewMessage(m *protogen.Message, c *cfg.Config) *Message {
// 	return &Message{Message: m, Config: c}
// }

// type Message struct {
// 	*protogen.Message
// 	Config *cfg.Config
// }

// func (m *Message) MessageRule() *authorize.Rule {
// 	return proto.GetExtension(m.Desc.Options(), authorize.E_Rule).(*authorize.Rule)
// }

// func (m *Message) Macros() ([]string, error) {
// 	if m.MessageRule() == nil {
// 		return nil, nil
// 	}
// 	envOpts := []cel.EnvOption{
// 		cel.Types(&authorize.AuthorizationContext{}),
// 		cel.DeclareContextProto(m.Desc),
// 		cel.Declarations(
// 			decls.NewVar("_ctx", decls.NewObjectType(string((&authorize.AuthorizationContext{}).ProtoReflect().Descriptor().FullName()))),
// 		),
// 	}
// 	for k := range m.Config.Globals.Functions {
// 		envOpts = append(envOpts, cel.Declarations(decls.NewFunction(k, decls.NewOverload(k, []*v1alpha1.Type{}, &v1alpha1.Type{TypeKind: &v1alpha1.Type_Primitive{Primitive: v1alpha1.Type_BOOL}}))))
// 	}
// 	env, err := cel.NewEnv(envOpts...)
// 	if err != nil {
// 		return nil, err
// 	}
// 	ast, issues := env.Compile(m.MessageRule().Expr)
// 	if issues != nil && issues.Err() != nil {
// 		return nil, issues.Err()
// 	}
// 	return findMacrosInAST(ast, m.Config.Globals.Functions), nil
// }

// func (m *Message) Validate() error {
// 	if m.MessageRule() == nil {
// 		return nil
// 	}
// 	mm := map[string]string{}
// 	if macros, err := m.Macros(); err != nil {
// 		return err
// 	} else {
// 		for _, macro := range macros {
// 			mm[macro] = m.Config.Globals.Functions[macro]
// 		}
// 	}
// 	envOpts, err := authorize.BuildEnvOptionsWithMacros([]cel.EnvOption{
// 		cel.Types(&authorize.AuthorizationContext{}),
// 		cel.DeclareContextProto(m.Desc),
// 		cel.Declarations(
// 			decls.NewVar("_ctx", decls.NewObjectType(string((&authorize.AuthorizationContext{}).ProtoReflect().Descriptor().FullName()))),
// 		),
// 	}, mm)
// 	if err != nil {
// 		return err
// 	}
// 	env, err := cel.NewEnv(envOpts...)
// 	if err != nil {
// 		return err
// 	}
// 	ast, issues := env.Compile(m.MessageRule().Expr)
// 	if issues != nil && issues.Err() != nil {
// 		return issues.Err()
// 	}
// 	switch ast.ResultType().TypeKind.(type) {
// 	case *v1alpha1.Type_Primitive:
// 		if v1alpha1.Type_PrimitiveType(ast.ResultType().GetPrimitive().Number()) != v1alpha1.Type_BOOL {
// 			return fmt.Errorf("result type not bool")
// 		}
// 	default:
// 		return fmt.Errorf("result type not bool")
// 	}
// 	_, err = env.Program(ast, cel.OptimizeRegex(interpreter.MatchesRegexOptimization))
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
