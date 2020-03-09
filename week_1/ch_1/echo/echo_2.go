package main

import (
	"fmt"
	"os"
)

func main() {

	s, sep := "", ""
	for _, arg := range os.Args[1:] {

		s += sep + arg
		sep = " "
	}
	fmt.Println(s)
}

// $ go run week_1/ch_1/echo/echo_2.go 2 3 5 7
// 2 3 5 7
