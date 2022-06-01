package authorize

import (
	"testing"

	"github.com/google/cel-go/common/operators"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	v1alpha1 "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

func TestFindMacrosExpr(t *testing.T) {
	m := map[string]string{
		"myFn": "1 == 1",
	}
	tests := []struct {
		Name   string
		Expr   *v1alpha1.Expr
		Result []string
	}{
		{
			Name: "None",
			Expr: &v1alpha1.Expr{
				ExprKind: &v1alpha1.Expr_CallExpr{
					CallExpr: &v1alpha1.Expr_Call{
						Function: operators.Equals,
						Args: []*v1alpha1.Expr{
							{
								ExprKind: &v1alpha1.Expr_ConstExpr{
									ConstExpr: &v1alpha1.Constant{
										ConstantKind: &v1alpha1.Constant_Int64Value{
											Int64Value: 1,
										},
									},
								},
							},
							{
								ExprKind: &v1alpha1.Expr_ConstExpr{
									ConstExpr: &v1alpha1.Constant{
										ConstantKind: &v1alpha1.Constant_Int64Value{
											Int64Value: 1,
										},
									},
								},
							},
						},
					},
				},
			},
			Result: []string{},
		},
		{
			Name: "Simple",
			Expr: &v1alpha1.Expr{
				ExprKind: &v1alpha1.Expr_CallExpr{
					CallExpr: &v1alpha1.Expr_Call{
						Function: "myFn",
						Args:     []*v1alpha1.Expr{},
					},
				},
			},
			Result: []string{"myFn"},
		},
		{
			Name: "Embed",
			Expr: &v1alpha1.Expr{
				ExprKind: &v1alpha1.Expr_CallExpr{
					CallExpr: &v1alpha1.Expr_Call{
						Function: operators.In,
						Args: []*v1alpha1.Expr{
							{
								ExprKind: &v1alpha1.Expr_ConstExpr{
									ConstExpr: &v1alpha1.Constant{
										ConstantKind: &v1alpha1.Constant_BoolValue{
											BoolValue: true,
										},
									},
								},
							},
							{
								ExprKind: &v1alpha1.Expr_CallExpr{
									CallExpr: &v1alpha1.Expr_Call{
										Function: "myFn",
										Args:     []*v1alpha1.Expr{},
									},
								},
							},
						},
					},
				},
			},
			Result: []string{"myFn"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			res := findMacrosExpr(tt.Expr, m)
			if !cmp.Equal(res, tt.Result, cmpopts.SortSlices(func(l, r string) bool {
				return l < r
			})) {
				t.Errorf("want %v, got %v", tt.Result, res)
			}
		})
	}
}
