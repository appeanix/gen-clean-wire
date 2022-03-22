package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"os"

	"golang.org/x/tools/go/packages"
)

type Meta struct {
	// User-defined flags
	Project string
	// The path to the package containing the stub to generate
	PackagePath string
	// The output path of the generated file
	OutputPath string

	// non-user defined or programatic fields
	Template    string
	FileName    string
	ORM         bool
	RpcPackage  string
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
		meta.FileName = "gen_rules.go"
		genRules()
	case "usecases", "UseCases", "useCases":
		meta.Template = usecaseTemplate
		meta.FileName = "gen_usecases.go"
		genUseCases()
	case "rpc", "Rpc", "RPC":
		meta.Template = rpcTemplate
		meta.FileName = "gen_rpcs.go"
		genRpc()
	}

	buildTemplate()
	writeFile()
}

func buildTemplate() {

	var err error
	tmpl, err = tmpl.Parse(meta.Template)
	failErr(err)

	tmpl.Execute(&genBuffer, meta)
	failErr(err)

}

func writeFile() {
	var outPath string
	if len(meta.OutputPath) > 0 {
		outPath = meta.OutputPath
	} else {
		outPath = "gen"
	}

	err := os.MkdirAll(outPath, os.ModePerm)
	failErr(err)

	f, err := os.Create(fmt.Sprintf("%s/%s", outPath, meta.FileName))
	failErr(err)
	defer f.Close()

	_, err = f.Write(genBuffer.Bytes())
	failErr(err)

	f.Sync()
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
