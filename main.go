package main

import (
	"fmt"
	_ "github.com/xuri/excelize/v2"
)

type Test struct {
	sum   int
	count int
}

func main() {
	bf := functionByHand()
	fmt.Println(bf.Moebius())
	fmt.Println(bf.Moebius().printANF())
	//nonDegenerateTest()
}

func pow2(n int) int {
	return 1 << n
}
