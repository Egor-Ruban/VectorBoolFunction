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
	bf, _ := newRandomVBF(3, 3)
	fmt.Println(bf)
	fmt.Println(bf.Moebius())
	fmt.Printf("%032b\n", bf.isCoordinatesDegenerate())
	fmt.Println(bf.isNonDegenerate())
}

func pow2(n int) int {
	return 1 << n
}
