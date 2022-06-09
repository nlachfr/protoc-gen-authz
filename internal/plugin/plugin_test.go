package plugin

import (
	"fmt"
	"testing"

	"github.com/Neakxs/protoc-gen-authz/authorize"
)

func TestLoadConfig(t *testing.T) {
	c := &authorize.FileRule{}
	if err := LoadConfig("../../testdata/config.yml", c); err != nil {
		t.Errorf("want nil, got %v", err)
	}
	fmt.Println(c)
}
