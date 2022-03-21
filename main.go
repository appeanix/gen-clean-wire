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
	ORM      bool

	Rules    []RuleMeta
	UseCases []UseCaseMeta
}

var meta Meta
var genBuffer bytes.Buffer
var tmpl = template.New("file")

func main() {

	command := os.Args[1]
	os.Args = os.Args[1:]

	flag.StringVar(&meta.Project, "project", "", "Go project or module name")
	flag.BoolVar(&meta.ORM, "orm", false, "Support DB ORM ")
	flag.Parse()

	switch command {
	case "rules":
		meta.Template = ruleTemplate
		genRules()
	case "usecases", "UseCases", "useCases":
		meta.Template = usecaseTemplate
		genUseCases()
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
