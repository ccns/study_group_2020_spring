package main

import "fmt"

func main() {
	//initialization
	s := make([]string, 3)
	fmt.Println("initial:", s)

	s[0] = "a"
	s[1] = "b"
	s[2] = "c"
	fmt.Println("set:", s)
	fmt.Println("get:", s[2])
	fmt.Println("len:", len(s))
	fmt.Println("cap:", cap(s))

	//append
	s = append(s, "d")
	s = append(s, "e", "f")
	fmt.Println("appd:", s)

	//copy
	c := make([]string, len(s))
	copy(c, s)
	fmt.Println("copy:", c)

	//slice
	l := s[2:5]
	fmt.Println("slc:", l)

	// declare
	t := []string{"g", "h", "i"}
	fmt.Println("dcl:", t)

	// 2-D
	twoD := make([][]int, 3)
	for i := 0; i < 3; i++ {
		innerLen := i + 1
		twoD[i] = make([]int, innerLen)
		for j := 0; j < innerLen; j++ {
			twoD[i][j] = i + j
		}
	}
	fmt.Println("2d: ", twoD)
}
