package main

import "fmt"

func main() {
	s := "CCNSccns"
	fmt.Println(len(s))     // "8"
	fmt.Println(s[0], s[4]) // "67 99" ('C' and 'c')
        //ASCII 有 32 個控制字元，同字母大小寫因此差 32
	// c := s[len(s)]          // panic: index out of range
	// fmt.Println(c)

	fmt.Println("Hello, " + s[:4]) // "Hello, ccns"
}
