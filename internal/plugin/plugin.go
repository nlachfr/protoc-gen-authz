package plugin

import (
	"flag"
	"os"

	"github.com/Neakxs/protoc-gen-authz/internal/cfg"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoregistry"
	"gopkg.in/yaml.v2"
)

const PluginName = "protoc-gen-go-authz"

var PluginVersion = "0.0.0"

var (
	config = flag.String("config", "", "global configuration file")
)

func Run() {
	protogen.Options{
		ParamFunc: flag.CommandLine.Set,
	}.Run(func(gen *protogen.Plugin) error {
		c := cfg.Config{}
		if config != nil {
			b, err := os.ReadFile(*config)
			if err != nil {
				return err
			}
			if err := yaml.Unmarshal(b, &c); err != nil {
				return err
			}
		}
		var files protoregistry.Files
		for _, file := range gen.Files {
			if err := files.RegisterFile(file.Desc); err != nil {
				return err
			}
		}
		for _, file := range gen.Files {
			if !file.Generate {
				continue
			}
			if err := NewFile(gen, file, &c).Generate(); err != nil {
				return err
			}
		}
		return nil
	})
}
