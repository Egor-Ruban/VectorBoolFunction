package main

import (
	"fmt"
	_ "github.com/xuri/excelize/v2"
)

func main() {
	vbf, _ := newRandomVBF(3, 2)
	fmt.Println(vbf)
	fmt.Println(vbf.WHT())
	d, nap := vbf.Affine()
	fmt.Println(d)
	fmt.Println(nap)
}

func pow2(n int) int {
	return 1 << n
}
