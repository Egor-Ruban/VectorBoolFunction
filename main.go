package main

import (
	"fmt"
	"time"
)

func main() {
	test2()
}

func test2() {

	// Код для измерения

	vbf, _ := newRandomVBF(29, 1)
	//fmt.Println(vbf)
	start := time.Now()
	vbf.WHT()
	duration := time.Since(start)
	// Отформатированная строка,
	// например, "2h3m0.5s" или "4.503μs"
	fmt.Println(duration)
	//fmt.Println(vbf.WHT())

}

func pow2(n int) int {
	return 1 << n
}
