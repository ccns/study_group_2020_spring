package main

import (
	"fmt"
	"os"
)

func main() {
	const (
		_ = 1 << (10 * iota)
		KiB
		MiB
		GiB
		TiB
		PiB
		EiB
		ZiB
		YiB
	)
	fmt.Println(KiB) //1024
	fmt.Println(MiB) //1048576
	fmt.Println(GiB) //1073741824
	fmt.Println(TiB) //1099511627776
	fmt.Println(PiB) //1125899906842624
	fmt.Println(EiB) //1152921504606846976
	// fmt.Println(ZiB)   //exceeds 1 << 64 from here
	// fmt.Println(YiB)
	
	// *******cool*********
	fmt.Println(YiB / ZiB) //1024
	os.Exit(0)
}
