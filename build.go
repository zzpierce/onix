package main

import (
	"fmt"
	"os"
)

func main() {
	// for idx, args := range os.Args {
	// 	fmt.Println("参数"+strconv.Itoa(idx)+":", args)
	// }
	args := os.Args
	if len(args) == 1 {
		fmt.Println("need file name")
		return
	}
	fileName := args[1]
	code, err := readFile(fileName)
	if err != nil {
		panic(err)
	}
	fmt.Println(code)
}

func readFile(fileName string) (string, error) {
	content, err := os.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
