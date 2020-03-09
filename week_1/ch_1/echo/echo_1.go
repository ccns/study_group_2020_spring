package main

import (
	"fmt"
	"os"
)

func main() {

	var s, sep string
	for i := 1; i < len(os.Args); i++ {

		s += sep + os.Args[i]
		sep = " "
	}
	fmt.Println(s)
}

// $ go run week_1/ch_1/echo/echo_1.go 2 3 5 7
// 2 3 5 7
