package plugin

import (
	"github.com/google/cel-go/cel"
	expr "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

func findMacrosInAST(ast *cel.Ast, m map[string]string) []string {
	s := findMacrosInExpr(ast.Expr(), m)
	mm := map[string]bool{}
	for _, i := range s {
		mm[i] = true
	}
	s = []string{}
	for k, _ := range mm {
		s = append(s, k)
	}
	return s
}

func findMacrosInExpr(e *expr.Expr, m map[string]string) []string {
	res := []string{}
	switch exp := e.ExprKind.(type) {
	case *expr.Expr_ConstExpr:
	case *expr.Expr_IdentExpr:
	case *expr.Expr_SelectExpr:
	case *expr.Expr_CallExpr:
		if _, ok := m[exp.CallExpr.Function]; ok {
			res = append(res, exp.CallExpr.Function)
		} else {
			for _, i := range exp.CallExpr.Args {
				res = append(res, findMacrosInExpr(i, m)...)
			}
		}
	case *expr.Expr_ListExpr:
		for _, i := range exp.ListExpr.Elements {
			res = append(res, findMacrosInExpr(i, m)...)
		}
	case *expr.Expr_StructExpr:
	case *expr.Expr_ComprehensionExpr:
	}
	return res
}
