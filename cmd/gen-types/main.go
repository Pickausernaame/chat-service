package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) != 4 {
		log.Fatalf("invalid args count: %d", len(os.Args)-1)
	}

	pkg, types, out := os.Args[1], strings.Split(os.Args[2], ","), os.Args[3]
	if err := run(pkg, types, out); err != nil {
		log.Fatal(err)
	}

	p, _ := os.Getwd()
	log.Printf("%v generated\n", filepath.Join(p, out))
}

func run(pkg string, types []string, outFile string) error {
	f, err := os.OpenFile(outFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o755)
	if err != nil {
		return err
	}
	return executor.Execute(f, struct {
		PackageName string
		TypeNames   []string
	}{PackageName: pkg, TypeNames: types})
}
