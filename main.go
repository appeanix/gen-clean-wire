package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"

	"golang.org/x/tools/go/packages"
)

type Meta struct {
	Template string
	Project  string
	// The path to the package containing the stub to generate
	PackagePath string
	// The output path of the generated file
	OutputPath string
	ORM        bool
	RpcPackage string

	Rules       []RuleMeta
	UseCases    []UseCaseMeta
	RpcServices []RpcServiceMeta
}

var meta Meta
var genBuffer bytes.Buffer
var tmpl = template.New("file")

func main() {

	command := os.Args[1]
	os.Args = os.Args[1:]

	flag.StringVar(&meta.Project, "project", "", "Go project or module name")
	flag.StringVar(&meta.PackagePath, "pkgPath", "", "Package path containing stub to generate")
	flag.StringVar(&meta.OutputPath, "outPath", "", "Output path that the generated files are stored")

	flag.Parse()

	switch command {
	case "rules":
		meta.Template = ruleTemplate
		genRules()
	case "usecases", "UseCases", "useCases":
		meta.Template = usecaseTemplate
		genUseCases()
	case "rpc", "Rpc", "RPC":
		meta.Template = rpcTemplate
		genRpc()
	}

	write()
	log.Printf("%s", genBuffer.String())
}

func write() {

	var err error
	if tmpl, err = tmpl.Parse(meta.Template); err != nil {
		fmt.Println(err)
	}
	tmpl.Execute(&genBuffer, meta)
}

func loadPackage(path string) *packages.Package {
	cfg := &packages.Config{Mode: packages.NeedTypes | packages.NeedImports}
	pkgs, err := packages.Load(cfg, path)
	if err != nil {
		failErr(fmt.Errorf("loading packages for inspection: %v", err))
	}
	if packages.PrintErrors(pkgs) > 0 {
		os.Exit(1)
	}

	return pkgs[0]
}

func failErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
