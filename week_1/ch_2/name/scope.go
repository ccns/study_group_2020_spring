package main

import (
	"fmt"
)

var _foo int = 0

func foo() {

	// unused, because of not accessable out of func
	// var _foo int = 1
}

func main() {

	var _bar int = 2
	fmt.Println("_foo: ", _foo)
	fmt.Println("_bar: ", _bar)
}

// $ go run week_1/ch_2/name/scope.go
// _foo:  0
// _bar:  2
