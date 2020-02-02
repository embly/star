package main

import (
	"log"
	"os"

	"github.com/embly/star"
	"github.com/embly/star/src"
)

func main() {
	star.AddPackages(src.Packages)
	if len(os.Args) > 1 {
		if err := star.RunScript(os.Args[1]); err != nil {
			log.Fatal(err)
		}
	}
	star.REPL()

}
