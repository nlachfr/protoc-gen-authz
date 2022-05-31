package authorize

import (
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter/functions"
)

type Overload struct {
	Name     string
	Function func(...ref.Val) ref.Val
}

func BuildRuntimeLibrary(config *FileRule, overloads ...*Overload) cel.Library {
	ols := []*functions.Overload{}
	if config != nil && config.Overloads != nil {
		for i := 0; i < len(overloads); i++ {
			o := overloads[i]
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
				ols = append(ols, overload)
			}
		}
	}
	return &library{envOpts: []cel.EnvOption{}, pgrOpts: []cel.ProgramOption{cel.Functions(ols...)}}
}

type library struct {
	envOpts []cel.EnvOption
	pgrOpts []cel.ProgramOption
}

func (l *library) CompileOptions() []cel.EnvOption     { return l.envOpts }
func (l *library) ProgramOptions() []cel.ProgramOption { return l.pgrOpts }
