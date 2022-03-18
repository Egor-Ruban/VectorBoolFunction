package main

import (
	"fmt"
)

func main() {
	test2()
}

func test2() {

	// Код для измерения

	vbf, _ := newRandomVBF(3, 1)
	fmt.Println(vbf)

	fmt.Println(vbf.WHT())

}

func pow2(n int) int {
	return 1 << n
}
