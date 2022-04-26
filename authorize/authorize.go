package authorize

import (
	"context"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common"
	"github.com/google/cel-go/parser"
	expr "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/proto"
)

func AuthorizationContextFromContext(ctx context.Context) *AuthorizationContext {
	res := &AuthorizationContext{
		Peer: &AuthorizationContext_Peer{
			Addr:     "",
			AuthInfo: "",
		},
		Metadata: make(map[string]*AuthorizationContext_MetadataValue),
	}
	if p, ok := peer.FromContext(ctx); ok {
		res.Peer.Addr = p.Addr.String()
		res.Peer.AuthInfo = p.AuthInfo.AuthType()
	}
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		for k, v := range md {
			res.Metadata[k] = &AuthorizationContext_MetadataValue{
				Values: v,
			}
		}
	}
	return res
}

func BuildProgramVars(ctx context.Context, message proto.Message) interface{} {
	res := map[string]interface{}{
		"_ctx": AuthorizationContextFromContext(ctx),
	}
	fields := message.ProtoReflect().Descriptor().Fields()
	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		res[field.TextName()] = message.ProtoReflect().Get(field)
	}
	return res
}

func BuildEnvOptionsWithMacros(opts []cel.EnvOption, m map[string]string) ([]cel.EnvOption, error) {
	env, err := cel.NewEnv(opts...)
	if err != nil {
		return nil, err
	}
	macros := []parser.Macro{}
	for k, v := range m {
		ast, issues := env.Compile(v)
		if issues != nil && issues.Err() != nil {
			return nil, issues.Err()
		}
		macros = append(macros, parser.NewGlobalMacro(k, 0, BuildMacroExpander(ast)))
	}
	return append(opts, cel.Macros(macros...)), nil
}

func BuildMacroExpander(ast *cel.Ast) parser.MacroExpander {
	return func(eh parser.ExprHelper, target *expr.Expr, args []*expr.Expr) (*expr.Expr, *common.Error) {
		return translateExpr(ast.Expr(), eh), nil
	}
}

func translateExpr(e *expr.Expr, eh parser.ExprHelper) *expr.Expr {
	switch exp := e.ExprKind.(type) {
	case *expr.Expr_ConstExpr:
		switch k := exp.ConstExpr.ConstantKind.(type) {
		case *expr.Constant_BoolValue:
			return eh.LiteralBool(k.BoolValue)
		case *expr.Constant_Int64Value:
			return eh.LiteralInt(k.Int64Value)
		case *expr.Constant_Uint64Value:
			return eh.LiteralUint(k.Uint64Value)
		case *expr.Constant_DoubleValue:
			return eh.LiteralDouble(k.DoubleValue)
		case *expr.Constant_StringValue:
			return eh.LiteralString(k.StringValue)
		case *expr.Constant_BytesValue:
			return eh.LiteralBytes(k.BytesValue)
		default:
			return e
		}
	case *expr.Expr_IdentExpr:
		return eh.Ident(exp.IdentExpr.GetName())
	case *expr.Expr_SelectExpr:
		return eh.Select(translateExpr(exp.SelectExpr.GetOperand(), eh), exp.SelectExpr.GetField())
	case *expr.Expr_CallExpr:
		args := []*expr.Expr{}
		for i := 0; i < len(exp.CallExpr.Args); i++ {
			args = append(args, translateExpr(exp.CallExpr.Args[i], eh))
		}
		return eh.GlobalCall(exp.CallExpr.GetFunction(), args...)
	case *expr.Expr_ListExpr:
		args := []*expr.Expr{}
		for i := 0; i < len(exp.ListExpr.GetElements()); i++ {
			args = append(args, translateExpr(exp.ListExpr.Elements[i], eh))
		}
		return eh.NewList(args...)
	case *expr.Expr_StructExpr:
		fieldInits := []*expr.Expr_CreateStruct_Entry{}
		for i := 0; i < len(exp.StructExpr.Entries); i++ {
			entry := exp.StructExpr.Entries[i]
			switch eexp := entry.KeyKind.(type) {
			case *expr.Expr_CreateStruct_Entry_FieldKey:
				fieldInits = append(fieldInits, eh.NewObjectFieldInit(eexp.FieldKey, entry.Value))
			case *expr.Expr_CreateStruct_Entry_MapKey:
				fieldInits = append(fieldInits, eh.NewMapEntry(eexp.MapKey, entry.Value))
			}
		}
		return eh.NewObject(exp.StructExpr.MessageName, fieldInits...)
	case *expr.Expr_ComprehensionExpr:
		return eh.Fold(
			exp.ComprehensionExpr.IterVar,
			translateExpr(exp.ComprehensionExpr.IterRange, eh),
			exp.ComprehensionExpr.AccuVar,
			translateExpr(exp.ComprehensionExpr.AccuInit, eh),
			translateExpr(exp.ComprehensionExpr.LoopCondition, eh),
			translateExpr(exp.ComprehensionExpr.LoopStep, eh),
			translateExpr(exp.ComprehensionExpr.Result, eh),
		)
	}
	return nil
}
