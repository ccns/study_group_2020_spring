package main

import "fmt"

func main() {

	var b bool

	b = 1       // error
	b = bool(1) // error

	b = (1 != 0) // correct
	fmt.Println(b)
}
