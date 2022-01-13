package main

import (
	"github.com/Neakxs/protoc-gen-authz/internal/plugin"
	"google.golang.org/protobuf/compiler/protogen"
)

func main() {
	protogen.Options{}.Run(plugin.Run)
}
