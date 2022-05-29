package authorize

import (
	"testing"

	"github.com/Neakxs/protoc-gen-authz/testdata"
)

func TestBuildAuthProgram(t *testing.T) {
	tests := []struct {
		Name    string
		Expr    string
		Config  *FileRule
		WantErr bool
	}{
		{
			Name:    "Unknown field",
			Expr:    `request.pong`,
			Config:  nil,
			WantErr: true,
		},
		{
			Name:    "Invalid return type",
			Expr:    `request.ping`,
			Config:  nil,
			WantErr: true,
		},
		{
			Name:    "OK",
			Expr:    `request.ping == "ping"`,
			Config:  nil,
			WantErr: false,
		},
		{
			Name: "OK (with bool macro)",
			Expr: `rule()`,
			Config: &FileRule{
				Globals: &FileRule_Globals{
					Functions: map[string]string{
						"rule": `request.ping == "ping"`,
					},
				},
			},
			WantErr: false,
		},
		{
			Name: "OK (with str macro)",
			Expr: `rule() == "ping"`,
			Config: &FileRule{
				Globals: &FileRule_Globals{
					Functions: map[string]string{
						"rule": `request.ping`,
					},
				},
			},
			WantErr: false,
		},
		{
			Name: "OK (array with str macro)",
			Expr: `"ping" in [rule()]`,
			Config: &FileRule{
				Globals: &FileRule_Globals{
					Functions: map[string]string{
						"rule": `request.ping`,
					},
				},
			},
			WantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			_, err := BuildAuthzProgram(tt.Expr, &testdata.PingRequest{}, tt.Config)
			if (tt.WantErr && err == nil) || (!tt.WantErr && err != nil) {
				t.Errorf("wantErr %v, got %v", tt.WantErr, err)
			}
		})
	}
}