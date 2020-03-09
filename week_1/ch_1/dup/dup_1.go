package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {

	counts := make(map[string]int)
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		counts[input.Text()]++
	}

	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d\t%s\n", n, line)
		}
	}
}

// $ go run week_1/ch_1/dup/dup_1.go < week_1/ch_1/dup/input_file
// 4       repeat 4 line
// 2       repeat 2 line
