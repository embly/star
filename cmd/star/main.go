package main

import (
	"os"

	"github.com/embly/star"
	"github.com/embly/star/src"
)

func main() {
	star.AddPackages(src.Packages)
	star.RunScript(os.Args[1])
}
