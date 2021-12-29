package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

func readFile(fileName string) (string, error) {
	content, err := os.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func main() {
	src, _ := readFile("../data/hello.go")

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		panic(err)
	}
	ast.Print(fset, f)
}
