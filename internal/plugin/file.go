package plugin

import (
	"github.com/Neakxs/protoc-gen-authz/authorize"
	"github.com/Neakxs/protoc-gen-authz/internal/template"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func NewFile(p *protogen.Plugin, f *protogen.File, c *authorize.FileRule) *File {
	g := p.NewGeneratedFile(f.GeneratedFilenamePrefix+".pb.authz.go", f.GoImportPath)
	cfg := &authorize.FileRule{}
	proto.Merge(cfg, c)
	fileRule := proto.GetExtension(f.Desc.Options(), authorize.E_File).(*authorize.FileRule)
	if fileRule != nil {
		proto.Merge(cfg, fileRule)
	}
	svcs := []*Service{}
	for i := 0; i < len(f.Services); i++ {
		svcs = append(svcs, NewService(f.Services[i], cfg, p.Files))
	}
	return &File{
		p:        p,
		g:        g,
		File:     f,
		Config:   cfg,
		Services: svcs,
	}
}

type File struct {
	p *protogen.Plugin
	g *protogen.GeneratedFile
	*protogen.File
	Config   *authorize.FileRule
	Services []*Service
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

func NewService(s *protogen.Service, c *authorize.FileRule, imports []*protogen.File) *Service {
	mths := []*Method{}
	for i := 0; i < len(s.Methods); i++ {
		mths = append(mths, NewMethod(s.Methods[i], c, imports))
	}
	return &Service{Service: s, Config: c, Methods: mths, Imports: imports}
}

type Service struct {
	*protogen.Service
	Imports []*protogen.File
	Config  *authorize.FileRule
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

func NewMethod(m *protogen.Method, c *authorize.FileRule, imports []*protogen.File) *Method {
	return &Method{Method: m, Config: c, Imports: imports}
}

type Method struct {
	*protogen.Method
	Imports []*protogen.File
	Config  *authorize.FileRule
}

func (m *Method) MethodRule() *authorize.MethodRule {
	var rule *authorize.MethodRule
	if m.Config != nil {
		if r, ok := m.Config.Rules[string(m.Desc.FullName())]; ok {
			rule = r
		}
	}
	if r, ok := proto.GetExtension(m.Desc.Options(), authorize.E_Method).(*authorize.MethodRule); ok && r != nil {
		rule = r
	}
	return rule
}

func (m *Method) Validate() error {
	rule := m.MethodRule()
	if rule == nil {
		return nil
	}
	imports := []protoreflect.FileDescriptor{}
	for i := 0; i < len(m.Imports); i++ {
		imports = append(imports, m.Imports[i].Desc)
	}
	if _, err := authorize.BuildAuthzProgramFromDesc(rule.GetExpr(), imports, m.Input.Desc, m.Config); err != nil {
		return err
	}
	return nil
}
