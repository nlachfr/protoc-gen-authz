package authorize

import (
	"testing"

	"github.com/Neakxs/protoc-gen-authz/testdata"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
)

type testOpt struct {
	fns  []*FunctionOverload
	vars []*VariableOverload
}

func (o *testOpt) GetFunctionOverloads() []*FunctionOverload { return o.fns }
func (o *testOpt) GetVariableOverloads() []*VariableOverload { return o.vars }

func TestBuildRuntimeLibrary(t *testing.T) {
	tests := []struct {
		Name    string
		Rule    string
		Config  *FileRule
		Options Options
		WantErr bool
	}{
		{
			Name: "Missing function overload",
			Rule: `myVariable == myFunction("")`,
			Config: &FileRule{
				Overloads: &FileRule_Overloads{
					Functions: map[string]*FileRule_Overloads_Function{
						"myFunction": {
							Args: []*FileRule_Overloads_Type{{
								Type: &FileRule_Overloads_Type_Primitive_{
									Primitive: FileRule_Overloads_Type_STRING,
								},
							}},
							Result: &FileRule_Overloads_Type{
								Type: &FileRule_Overloads_Type_Primitive_{
									Primitive: FileRule_Overloads_Type_STRING,
								},
							},
						},
					},
					Variables: map[string]*FileRule_Overloads_Type{
						"myVariable": {
							Type: &FileRule_Overloads_Type_Primitive_{
								Primitive: FileRule_Overloads_Type_STRING,
							},
						},
					},
				},
			},
			Options: &testOpt{
				vars: []*VariableOverload{{
					Name:  "myVariable",
					Value: "ok",
				}},
			},
			WantErr: true,
		},
		{
			Name: "Missing variable overload",
			Rule: `myVariable == myFunction("")`,
			Config: &FileRule{
				Overloads: &FileRule_Overloads{
					Functions: map[string]*FileRule_Overloads_Function{
						"myFunction": {
							Args: []*FileRule_Overloads_Type{{
								Type: &FileRule_Overloads_Type_Primitive_{
									Primitive: FileRule_Overloads_Type_STRING,
								},
							}},
							Result: &FileRule_Overloads_Type{
								Type: &FileRule_Overloads_Type_Primitive_{
									Primitive: FileRule_Overloads_Type_STRING,
								},
							},
						},
					},
					Variables: map[string]*FileRule_Overloads_Type{
						"myVariable": {
							Type: &FileRule_Overloads_Type_Primitive_{
								Primitive: FileRule_Overloads_Type_STRING,
							},
						},
					},
				},
			},
			Options: &testOpt{
				fns: []*FunctionOverload{{
					Name: "myFunction",
					Function: func(v ...ref.Val) ref.Val {
						return types.String("ok")
					},
				}},
			},
			WantErr: true,
		},
		{
			Name: "OK (1 arg)",
			Rule: `myVariable == myFunction("")`,
			Config: &FileRule{
				Overloads: &FileRule_Overloads{
					Functions: map[string]*FileRule_Overloads_Function{
						"myFunction": {
							Args: []*FileRule_Overloads_Type{{
								Type: &FileRule_Overloads_Type_Primitive_{
									Primitive: FileRule_Overloads_Type_STRING,
								},
							}},
							Result: &FileRule_Overloads_Type{
								Type: &FileRule_Overloads_Type_Primitive_{
									Primitive: FileRule_Overloads_Type_STRING,
								},
							},
						},
					},
					Variables: map[string]*FileRule_Overloads_Type{
						"myVariable": {
							Type: &FileRule_Overloads_Type_Primitive_{
								Primitive: FileRule_Overloads_Type_STRING,
							},
						},
					},
				},
			},
			Options: &testOpt{
				fns: []*FunctionOverload{{
					Name: "myFunction",
					Function: func(v ...ref.Val) ref.Val {
						return types.String("ok")
					},
				}},
				vars: []*VariableOverload{{
					Name:  "myVariable",
					Value: "ok",
				}},
			},
			WantErr: false,
		},
		{
			Name: "OK (2 args)",
			Rule: `myVariable == myFunction("", "")`,
			Config: &FileRule{
				Overloads: &FileRule_Overloads{
					Functions: map[string]*FileRule_Overloads_Function{
						"myFunction": {
							Args: []*FileRule_Overloads_Type{{
								Type: &FileRule_Overloads_Type_Primitive_{
									Primitive: FileRule_Overloads_Type_STRING,
								},
							}, {
								Type: &FileRule_Overloads_Type_Primitive_{
									Primitive: FileRule_Overloads_Type_STRING,
								},
							}},
							Result: &FileRule_Overloads_Type{
								Type: &FileRule_Overloads_Type_Primitive_{
									Primitive: FileRule_Overloads_Type_STRING,
								},
							},
						},
					},
					Variables: map[string]*FileRule_Overloads_Type{
						"myVariable": {
							Type: &FileRule_Overloads_Type_Primitive_{
								Primitive: FileRule_Overloads_Type_STRING,
							},
						},
					},
				},
			},
			Options: &testOpt{
				fns: []*FunctionOverload{{
					Name: "myFunction",
					Function: func(v ...ref.Val) ref.Val {
						return types.String("ok")
					},
				}},
				vars: []*VariableOverload{{
					Name:  "myVariable",
					Value: "ok",
				}},
			},
			WantErr: false,
		},
		{
			Name: "OK (any args)",
			Rule: "myVariable == myFunction()",
			Config: &FileRule{
				Overloads: &FileRule_Overloads{
					Functions: map[string]*FileRule_Overloads_Function{
						"myFunction": {
							Args: []*FileRule_Overloads_Type{},
							Result: &FileRule_Overloads_Type{
								Type: &FileRule_Overloads_Type_Primitive_{
									Primitive: FileRule_Overloads_Type_STRING,
								},
							},
						},
					},
					Variables: map[string]*FileRule_Overloads_Type{
						"myVariable": {
							Type: &FileRule_Overloads_Type_Primitive_{
								Primitive: FileRule_Overloads_Type_STRING,
							},
						},
					},
				},
			},
			Options: &testOpt{
				fns: []*FunctionOverload{{
					Name: "myFunction",
					Function: func(v ...ref.Val) ref.Val {
						return types.String("ok")
					},
				}},
				vars: []*VariableOverload{{
					Name:  "myVariable",
					Value: "ok",
				}},
			},
			WantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			pgr, err := BuildAuthzProgram(tt.Rule, &testdata.PingRequest{}, tt.Config, BuildRuntimeLibrary(tt.Config, tt.Options))
			if err != nil {
				if !tt.WantErr {
					t.Errorf("wantErr %v, got %v", tt.WantErr, err)
				}
			}
			_, _, err = pgr.Eval(map[string]interface{}{})
			if (err == nil && tt.WantErr) || (!tt.WantErr && err != nil) {
				t.Errorf("wantErr %v, got %v", tt.WantErr, err)
			}
		})
	}
}
