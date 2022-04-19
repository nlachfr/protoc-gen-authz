package plugin

// import (
// 	"fmt"
// 	"strings"

// 	"github.com/google/cel-go/cel"
// 	"github.com/google/cel-go/common/operators"
// 	expr "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
// 	"google.golang.org/protobuf/reflect/protoreflect"
// )

// func TranspileAST(ast *cel.Ast, messageDesc protoreflect.MessageDescriptor) (string, error) {
// 	fields := messageDesc.Fields()
// 	for i := 0; i < fields.Len(); i++ {
// 		fieldDesc := fields.Get(i)

// 	}

// 	return transpileExpr(ast, ast.Expr())
// }

// func transpileConstExpr(ast *cel.Ast, exp *expr.Expr) (string, error) {
// 	return "", nil
// }

// func transpileIdentExpr(ast *cel.Ast, exp *expr.Expr) (string, error) {
// 	identExp := exp.GetIdentExpr()
// 	if identExp == nil {
// 		return "", fmt.Errorf("ident expr is nil")
// 	}
// 	return "", nil
// }

// func transpileCallExpr(ast *cel.Ast, exp *expr.Expr) (string, error) {
// 	callExp := exp.GetCallExpr()
// 	if callExp == nil {
// 		return "", fmt.Errorf("call expr is nil")
// 	}
// 	args := []string{}
// 	for i := 0; i < len(callExp.Args); i++ {
// 		if res, err := transpileExpr(ast, callExp.Args[i]); err != nil {
// 			return "", err
// 		} else {
// 			args = append(args, res)
// 		}
// 	}
// 	switch callExp.Function {
// 	case operators.Conditional:
// 	case operators.LogicalAnd:
// 		return fmt.Sprintf("(%s)", strings.Join(args, " && ")), nil
// 	case operators.LogicalOr:
// 		return fmt.Sprintf("(%s)", strings.Join(args, " || ")), nil
// 	case operators.LogicalNot:
// 		return "!" + args[0], nil
// 	case operators.Equals:
// 		return fmt.Sprintf("(%s)", strings.Join(args, " == ")), nil
// 	case operators.NotEquals:
// 		return fmt.Sprintf("(%s)", strings.Join(args, " != ")), nil
// 	case operators.Less:
// 		return fmt.Sprintf("(%s)", strings.Join(args, " < ")), nil
// 	case operators.LessEquals:
// 		return fmt.Sprintf("(%s)", strings.Join(args, " <= ")), nil
// 	case operators.Greater:
// 		return fmt.Sprintf("(%s)", strings.Join(args, " > ")), nil
// 	case operators.GreaterEquals:
// 		return fmt.Sprintf("(%s)", strings.Join(args, " >= ")), nil
// 	case operators.Add:
// 		return fmt.Sprintf("(%s)", strings.Join(args, " + ")), nil
// 	case operators.Subtract:
// 		return fmt.Sprintf("(%s)", strings.Join(args, " - ")), nil
// 	case operators.Multiply:
// 		return fmt.Sprintf("(%s)", strings.Join(args, " * ")), nil
// 	case operators.Divide:
// 		return fmt.Sprintf("(%s)", strings.Join(args, " / ")), nil
// 	case operators.Modulo:
// 		return fmt.Sprintf("(%s)", strings.Join(args, " % ")), nil
// 	case operators.Negate:
// 		return fmt.Sprintf("(%s)", strings.Join(args, " - ")), nil
// 	case operators.Index:
// 		return args[0] + "[" + args[1] + "]", nil
// 	}
// 	return "", nil
// }

// func transpileExpr(ast *cel.Ast, exp *expr.Expr) (string, error) {
// 	switch exp.ExprKind.(type) {
// 	case *expr.Expr_ConstExpr:
// 		return transpileConstExpr(ast, exp)
// 	case *expr.Expr_IdentExpr:
// 		return transpileIdentExpr(ast, exp)
// 	// case *expr.Expr_SelectExpr:
// 	case *expr.Expr_CallExpr:
// 		return transpileCallExpr(ast, exp)
// 	// case *expr.Expr_ListExpr:
// 	// case *expr.Expr_StructExpr:
// 	// case *expr.Expr_ComprehensionExpr:
// 	default:
// 		return "", fmt.Errorf("%t: expr kind not supported", exp.ExprKind)
// 	}
// }
