package main

import (
	"fmt"
	"go/build"
	"go/doc"
	"go/parser"
	"go/token"
	"log"
	"os"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	pkg, err := build.Import("io/ioutil", wd, build.ImportComment)
	if err != nil {
		log.Fatal(err)
	}

	fs := token.NewFileSet()
	include := func(info os.FileInfo) bool {
		for _, name := range pkg.GoFiles {
			if name == info.Name() {
				return true
			}
		}
		for _, name := range pkg.CgoFiles {
			if name == info.Name() {
				return true
			}
		}
		return false
	}
	pkgs, err := parser.ParseDir(fs, pkg.Dir, include, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}
	astPkg := pkgs[pkg.Name]
	mode := doc.AllDecls
	docPkg := doc.New(astPkg, pkg.ImportPath, mode)

	fmt.Println(docPkg.Funcs[0].Decl.Type.Params.List[0].Type)
}
