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
	bf := specialFunction()
	fmt.Println(bf.autoCorrelation())
	fmt.Println(bf.maxPC())
}

func pow2(n int) int {
	return 1 << n
}

func Abs(n int) int {
	if n < 0 {
		return -n
	} else {
		return n
	}
}
