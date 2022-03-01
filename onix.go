package main

import (
	"fmt"
	"onix/lex/ast"
	"os"
	"strconv"
)

func main() {
	for idx, args := range os.Args {
		fmt.Println("参数"+strconv.Itoa(idx)+":", args)
	}
	args := os.Args
	// if len(args) == 1 {
	// 	fmt.Println("need file name")
	// 	return
	// }
	fileName := "playground2/main.go"
	if len(args) > 1 {
		fileName = args[1]
	}

	atb := ast.AstTreeBuilder{}
	atb.Init(fileName)
	f, err := atb.BuildFile()
	if err != nil {
		panic(err)
	}
	fmt.Println(f)
}

func readFile(fileName string) (string, error) {
	content, err := os.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// func main() {
// 	s := "你来hello"
// 	for i := 0; i < len(s); i++ {
// 		fmt.Println(s[i])
// 	}
// }
