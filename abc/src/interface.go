package main

import "github.com/golang/protobuf/protoc-gen-go/generator"

type grpcPlugin struct {
	*generator.Generator
}

func (p *grpcPlugin) Name() string {
	return "grpc"
}

func (p *grpcPlugin) Init(g *generator.Generator) {
	p.Generator = g
}

func (p *grpcPlugin) GenerateImports(file *generator.FileDescriptor) {
	if len(file.Service) == 0 {
		return
	}
	p.P(`import "google.golang.org/grpc"`)
}

func main() {

}
