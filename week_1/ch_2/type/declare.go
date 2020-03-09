package main

import (
	"fmt"
)

type price float64
type height float64
type house struct {
	Price  price
	Height height
}

const (
	house_price  price  = 4000000.0
	house_height height = 35.5
)

func main() {

	// type mismatch may cause :
	// invalid operation: house_height < house_price (mismatched types height and price)
	// fmt.Println(house_height < house_price)

	fmt.Println(price(house_height) < house_price)
}

// $ go run week_1/ch_2/type/declare.go
// true
