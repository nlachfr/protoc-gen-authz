package authorize

import (
	"fmt"

	"github.com/Neakxs/protoc-gen-authz/cfg"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/interpreter"
	"github.com/google/cel-go/parser"
	v1alpha1 "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func BuildAuthzProgram(expr string, req interface{}, config *cfg.Config) (cel.Program, error) {
	var reqDesc protoreflect.MessageDescriptor
	var reqOpt cel.EnvOption
	if r, ok := req.(proto.Message); ok {
		reqOpt = cel.Types(r)
		reqDesc = r.ProtoReflect().Descriptor()
	} else if r, ok := req.(protoreflect.MessageDescriptor); ok {
		reqOpt = cel.TypeDescs(r.ParentFile())
		reqDesc = r
	} else {
		return nil, fmt.Errorf("invalid req")
	}
	authzCtx := &AuthorizationContext{}
	envOpts := []cel.EnvOption{
		reqOpt,
		cel.Types(authzCtx),
		cel.Declarations(
			decls.NewVar(
				"context",
				decls.NewObjectType(string(authzCtx.ProtoReflect().Descriptor().FullName()))),
			decls.NewVar(
				"request",
				decls.NewObjectType(string(reqDesc.FullName())),
			),
		),
	}
	macros := []parser.Macro{}
	if rawMacros, err := findMacros(config, envOpts, expr); err != nil {
		return nil, err
	} else {
		env, err := cel.NewEnv(envOpts...)
		if err != nil {
			return nil, err
		}
		for _, macro := range rawMacros {
			ast, issues := env.Compile(config.Globals.Functions[macro])
			if issues != nil && issues.Err() != nil {
				return nil, issues.Err()
			}
			macros = append(macros, parser.NewGlobalMacro(macro, 0, BuildMacroExpander(ast)))
		}
	}
	envOpts = append(envOpts, cel.Macros(macros...))
	env, err := cel.NewEnv(envOpts...)
	if err != nil {
		return nil, err
	}
	ast, issues := env.Compile(expr)
	if issues != nil && issues.Err() != nil {
		return nil, issues.Err()
	}
	switch ast.ResultType().TypeKind.(type) {
	case *v1alpha1.Type_Primitive:
		if v1alpha1.Type_PrimitiveType(ast.ResultType().GetPrimitive().Number()) != v1alpha1.Type_BOOL {
			return nil, fmt.Errorf("result type not bool")
		}
	default:
		return nil, fmt.Errorf("result type not bool")
	}
	pgr, err := env.Program(ast, cel.OptimizeRegex(interpreter.MatchesRegexOptimization))
	if err != nil {
		return nil, err
	}
	return pgr, nil
}

func findMacros(config *cfg.Config, opts []cel.EnvOption, expr string) ([]string, error) {
	envOpts := opts
	for k := range config.Globals.Functions {
		envOpts = append(envOpts, cel.Declarations(decls.NewFunction(k, decls.NewOverload(k, []*v1alpha1.Type{}, &v1alpha1.Type{TypeKind: &v1alpha1.Type_Dyn{}}))))
	}
	env, err := cel.NewEnv(envOpts...)
	if err != nil {
		return nil, err
	}
	ast, issues := env.Compile(expr)
	if issues != nil && issues.Err() != nil {
		return nil, issues.Err()
	}
	return findMacrosInAST(ast, config.Globals.Functions), nil
}

func findMacrosInAST(ast *cel.Ast, m map[string]string) []string {
	s := findMacrosInExpr(ast.Expr(), m)
	mm := map[string]bool{}
	for _, i := range s {
		mm[i] = true
	}
	s = []string{}
	for k := range mm {
		s = append(s, k)
	}
	return s
}

func findMacrosInExpr(e *v1alpha1.Expr, m map[string]string) []string {
	res := []string{}
	switch exp := e.ExprKind.(type) {
	case *v1alpha1.Expr_ConstExpr:
	case *v1alpha1.Expr_IdentExpr:
	case *v1alpha1.Expr_SelectExpr:
	case *v1alpha1.Expr_CallExpr:
		if _, ok := m[exp.CallExpr.Function]; ok {
			res = append(res, exp.CallExpr.Function)
		} else {
			for _, i := range exp.CallExpr.Args {
				res = append(res, findMacrosInExpr(i, m)...)
			}
		}
	case *v1alpha1.Expr_ListExpr:
		for _, i := range exp.ListExpr.Elements {
			res = append(res, findMacrosInExpr(i, m)...)
		}
	case *v1alpha1.Expr_StructExpr:
	case *v1alpha1.Expr_ComprehensionExpr:
	}
	return res
}
