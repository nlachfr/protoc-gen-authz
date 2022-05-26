package plugin

import (
	"github.com/Neakxs/protoc-gen-authz/authorize"
	"github.com/Neakxs/protoc-gen-authz/internal/template"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

func NewFile(p *protogen.Plugin, f *protogen.File, c *authorize.FileRule) *File {
	g := p.NewGeneratedFile(f.GeneratedFilenamePrefix+".pb.authz.go", f.GoImportPath)
	svcs := []*Service{}
	for i := 0; i < len(f.Services); i++ {
		svcs = append(svcs, NewService(f.Services[i], c))
	}
	cfg := &authorize.FileRule{}
	proto.Merge(cfg, c)
	fileRule := proto.GetExtension(f.Desc.Options(), authorize.E_File).(*authorize.FileRule)
	if fileRule != nil {
		proto.Merge(cfg, fileRule)
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

func NewService(s *protogen.Service, c *authorize.FileRule) *Service {
	mths := []*Method{}
	for i := 0; i < len(s.Methods); i++ {
		mths = append(mths, NewMethod(s.Methods[i], c))
	}
	return &Service{Service: s, Config: c, Methods: mths}
}

type Service struct {
	*protogen.Service
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

func NewMethod(m *protogen.Method, c *authorize.FileRule) *Method {
	return &Method{Method: m, Config: c}
}

type Method struct {
	*protogen.Method
	Config *authorize.FileRule
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
