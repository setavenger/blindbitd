package main

import (
	"flag"
	"github.com/setavenger/blindbit-cli/cmd"
	"github.com/spf13/cobra/doc"
	"log"
	"os"
)

var outDir string

func init() {
	flag.StringVar(&outDir, "out", "./docs", "set output directory")
	flag.Parse()
}

func main() {
	err := os.MkdirAll(outDir, 0755)
	if err != nil {
		log.Fatal(err)
		return
	}
	err = doc.GenMarkdownTree(cmd.RootCmd, outDir)
	if err != nil {
		log.Fatal(err)
	}
}
