package main

import (
	"fmt"
	"math"
	"strconv"
)

func main() {
	fmt.Println("hello world" + strconv.FormatFloat(math.Abs(100), 'f', -1, 64))
}

func add(a, b int) (c int) {
	c += 1
	return a + b
}
