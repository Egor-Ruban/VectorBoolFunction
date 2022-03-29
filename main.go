package main

import (
	_ "github.com/xuri/excelize/v2"
)

type Test struct {
	sum   int
	count int
}

func main() {
	nonDegenerateTest()
}

func pow2(n int) int {
	return 1 << n
}
