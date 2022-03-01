package ast

import "fmt"

func expectError(expect, actual string) error {
	return fmt.Errorf("expect '%s', found '%s'", expect, actual)
}

func unexpectError(actual string) error {
	return fmt.Errorf("unexpect identifier: %s", actual)
}
