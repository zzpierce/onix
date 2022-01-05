package main

import (
	"fmt"
	"onix/lex"
	"os"
	"strconv"
)

func main() {
	for idx, args := range os.Args {
		fmt.Println("参数"+strconv.Itoa(idx)+":", args)
	}
	args := os.Args
	if len(args) == 1 {
		fmt.Println("need file name")
		return
	}
	fileName := args[1]
	sc, err := lex.InitScan(fileName)
	if err != nil {
		panic(err)
	}
	for {
		lit, err := sc.Literal()
		if err != nil {
			panic(err)
		}
		if lit == "" {
			break
		}
		fmt.Println(lit)
	}
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
