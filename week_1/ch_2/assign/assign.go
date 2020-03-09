package main

import (
	"fmt"
)

func multipleReturn() (string, int) {

	return "first", 0
}

func main() {

	_, value := multipleReturn()
	value++
	fmt.Println("value: ", value)

	medals := []string{"gold", "silver", "bronze"}
	medals[0] = "Platinum"
	fmt.Println("medals: ", medals)
}

// $ go run week_1/ch_2/assign/assign.go
// value:  1
// medals:  [Platinum silver bronze]
