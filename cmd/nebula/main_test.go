package main

import (
	"flag"
	"fmt"
	"go/build"
	"go/doc"
	"go/parser"
	"go/token"
	"log"
	"os"
	"testing"
)

func TestPackage(t *testing.T) {
	dirsInit()
	if err := do(os.Stdout, flag.CommandLine, []string{"io/ioutil"}); err != nil {
		panic(err)
	}
}

func TestBasic(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	pkg, err := build.Import("io/ioutil", wd, build.ImportComment)
	if err != nil {
		t.Fatal(err)
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
