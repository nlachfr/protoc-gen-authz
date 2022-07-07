package authorize

import (
	"fmt"
	"net/http"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter"
	"github.com/google/cel-go/interpreter/functions"
	"github.com/google/cel-go/parser"
	v1alpha1 "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func BuildAuthzProgramFromDesc(expr string, imports []protoreflect.FileDescriptor, msgDesc protoreflect.MessageDescriptor, config *FileRule, libs ...cel.Library) (cel.Program, error) {
	envOpts := []cel.EnvOption{
		cel.TypeDescs(msgDesc.Parent()),
	}
	for i := 0; i < len(imports); i++ {
		envOpts = append(envOpts, cel.TypeDescs(imports[i]))
	}
	for i := 0; i < len(libs); i++ {
		envOpts = append(envOpts, cel.Lib(libs[i]))
	}
	return buildAuthzProgram(expr, msgDesc, config, envOpts...)
}

func BuildAuthzProgram(expr string, msg proto.Message, config *FileRule, libs ...cel.Library) (cel.Program, error) {
	envOpts := []cel.EnvOption{
		cel.Types(msg),
	}
	for i := 0; i < len(libs); i++ {
		envOpts = append(envOpts, cel.Lib(libs[i]))
	}
	return buildAuthzProgram(expr, msg.ProtoReflect().Descriptor(), config, envOpts...)
}

func buildAuthzProgram(expr string, desc protoreflect.MessageDescriptor, config *FileRule, envOpts ...cel.EnvOption) (cel.Program, error) {
	envOpts = append(envOpts,
		cel.Lib(&library{
			envOpts: []cel.EnvOption{
				cel.Declarations(
					decls.NewFunction(
						"get",
						decls.NewInstanceOverload(
							"get",
							[]*v1alpha1.Type{
								decls.NewMapType(
									decls.String,
									decls.NewListType(decls.String),
								),
								decls.String,
							}, decls.String,
						),
					),
					decls.NewFunction(
						"values",
						decls.NewInstanceOverload(
							"values",
							[]*v1alpha1.Type{
								decls.NewMapType(
									decls.String,
									decls.NewListType(decls.String),
								),
								decls.String,
							}, decls.NewListType(decls.String),
						),
					),
				),
			},
			pgrOpts: []cel.ProgramOption{
				cel.Functions(&functions.Overload{
					Operator: "get",
					Binary: func(lhs, rhs ref.Val) ref.Val {
						var h http.Header
						switch m := lhs.Value().(type) {
						case map[string][]string:
							h = http.Header(m)
						case http.Header:
							h = m
						default:
							return types.String("")
						}
						if s, ok := rhs.Value().(string); ok {
							return types.String(h.Get(s))
						}
						return types.String("")
					},
				}, &functions.Overload{
					Operator: "values",
					Binary: func(lhs, rhs ref.Val) ref.Val {
						var h http.Header
						switch m := lhs.Value().(type) {
						case map[string][]string:
							h = http.Header(m)
						case http.Header:
							h = m
						default:
							return types.NewStringList(nil, []string{})
						}
						if s, ok := rhs.Value().(string); ok {
							return types.NewStringList(TypeAdapterFunc(func(value interface{}) ref.Val { return types.String(value.(string)) }), h.Values(s))
						}
						return types.String("")
					},
				}),
			},
		}),
		cel.Declarations(
			decls.NewVar(
				"headers",
				decls.NewMapType(
					decls.String,
					decls.NewListType(decls.String),
				),
			),
			decls.NewVar(
				"request",
				decls.NewObjectType(string(desc.FullName())),
			),
		),
		buildGlobalConstantsOption(config),
		buildOverloadFunctionsOption(config),
		buildOverloadVariablesOption(config),
	)
	macros := []parser.Macro{}
	if rawMacros, err := FindMacros(config, envOpts, expr); err != nil {
		return nil, err
	} else {
		env, err := cel.NewEnv(envOpts...)
		if err != nil {
			return nil, fmt.Errorf("new env error: %w", err)
		}
		for _, macro := range rawMacros {
			ast, issues := env.Compile(config.Globals.Functions[macro])
			if issues != nil && issues.Err() != nil {
				return nil, fmt.Errorf("macro error: %w", issues.Err())
			}
			macros = append(macros, parser.NewGlobalMacro(macro, 0, BuildMacroExpander(ast)))
		}
	}
	envOpts = append(envOpts, cel.Macros(macros...))
	env, err := cel.NewEnv(envOpts...)
	if err != nil {
		return nil, fmt.Errorf("new env error: %w", err)
	}
	ast, issues := env.Compile(expr)
	if issues != nil && issues.Err() != nil {
		return nil, fmt.Errorf("compile error: %w", issues.Err())
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
		return nil, fmt.Errorf("program error: %w", err)
	}
	return pgr, nil
}

func buildGlobalConstantsOption(config *FileRule) cel.EnvOption {
	constantDecls := []*v1alpha1.Decl{}
	if config != nil && config.Globals != nil {
		for k, v := range config.Globals.Constants {
			constantDecls = append(constantDecls, decls.NewConst(
				k,
				decls.String,
				&v1alpha1.Constant{ConstantKind: &v1alpha1.Constant_StringValue{StringValue: v}},
			))
		}
	}
	return cel.Declarations(constantDecls...)
}

func buildOverloadFunctionsOption(config *FileRule) cel.EnvOption {
	functionDecls := []*v1alpha1.Decl{}
	if config != nil && config.Overloads != nil {
		for name, v := range config.Overloads.Functions {
			args := []*v1alpha1.Type{}
			overload := name
			for i := 0; i < len(v.Args); i++ {
				args = append(args, TypeFromOverloadType(v.Args[i]))
			}
			functionDecls = append(functionDecls, decls.NewFunction(
				name, decls.NewOverload(
					overload,
					args,
					TypeFromOverloadType(v.Result),
				),
			))
		}
	}
	return cel.Declarations(functionDecls...)
}

func buildOverloadVariablesOption(config *FileRule) cel.EnvOption {
	variableDecls := []*v1alpha1.Decl{}
	if config != nil && config.Overloads != nil {
		for k, v := range config.Overloads.Variables {
			variableDecls = append(variableDecls, decls.NewVar(k, TypeFromOverloadType(v)))
		}
	}
	return cel.Declarations(variableDecls...)
}
