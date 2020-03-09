package main

import (
	"fmt"
)

func incr(p *int) int {
	*p++
	return *p
}

func main() {

	v := 1
	incr(&v)
	fmt.Println(incr(&v))
}

// $ go run week_1/ch_2/pointer/order.go
// 3
