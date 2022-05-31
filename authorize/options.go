package authorize

import (
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter/functions"
)

type Options interface {
	GetFunctionOverloads() []*FunctionOverload
	GetVariableOverloads() []*VariableOverload
}

type FunctionOverload struct {
	Name     string
	Function func(...ref.Val) ref.Val
}

type VariableOverload struct {
	Name  string
	Value interface{}
}

func BuildRuntimeLibrary(config *FileRule, opts ...Options) cel.Library {
	fns := []*functions.Overload{}
	vars := map[string]interface{}{}
	if config != nil {
		if config.Overloads != nil {
			fnOverloads := []*FunctionOverload{}
			for i := 0; i < len(opts); i++ {
				fnOverloads = append(fnOverloads, opts[i].GetFunctionOverloads()...)
				vs := opts[i].GetVariableOverloads()
				for j := 0; j < len(vs); j++ {
					vars[vs[j].Name] = vs[j].Value
				}
			}
			for i := 0; i < len(fnOverloads); i++ {
				o := fnOverloads[i]
				if fnOverload, ok := config.Overloads.Functions[o.Name]; ok {
					overload := &functions.Overload{
						Operator: o.Name,
					}
					switch len(fnOverload.Args) {
					case 1:
						overload.Unary = func(value ref.Val) ref.Val {
							return o.Function(value)
						}
					case 2:
						overload.Binary = func(lhs, rhs ref.Val) ref.Val {
							return o.Function(lhs, rhs)
						}
					default:
						overload.Function = func(values ...ref.Val) ref.Val {
							return o.Function(values...)
						}
					}
					fns = append(fns, overload)
				}
			}
		}
	}
	return &library{envOpts: []cel.EnvOption{}, pgrOpts: []cel.ProgramOption{cel.Functions(fns...), cel.Globals(vars)}}
}

type library struct {
	envOpts []cel.EnvOption
	pgrOpts []cel.ProgramOption
}

func (l *library) CompileOptions() []cel.EnvOption     { return l.envOpts }
func (l *library) ProgramOptions() []cel.ProgramOption { return l.pgrOpts }
