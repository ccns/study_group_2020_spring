package main

import (
	"fmt"
	"reflect"
)

var foo, bar, mur = 0.2, true, "Hi"

func declare() {

	unique := 0
	unique, sub := 1, 2
	unique, sub = sub, unique

	// re-declare error
	// unique := 3

	fmt.Println("declare: ", unique, sub)
}

func main() {

	declare()
	fmt.Println(foo, bar, mur)
	fmt.Println(reflect.TypeOf(foo), reflect.TypeOf(bar), reflect.TypeOf(mur))
}

// $ go run week_1/ch_2/variable/declare.go
// declare:  2 1
// 0.2 true Hi
// float64 bool string
