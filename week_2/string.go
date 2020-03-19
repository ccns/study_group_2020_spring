package main

import "fmt"

func main() {
	s := "CCNSccns"
	fmt.Println(len(s))     // "8"
	fmt.Println(s[0], s[4]) // "67 99" ('C' and 'c')
	//ASCII有32個控制字元同字母大小寫因此差32
	// c := s[len(s)]          // panic: index out of range
	// fmt.Println(c)

	fmt.Println("Hello, " + s[:4]) // "Hello, ccns"
}
