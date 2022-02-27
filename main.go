package main

import "fmt"

func main() {
	test2()
}

func test2() {
	vbf, _ := newRandomVBF(3, 2)
	fmt.Println(vbf)
	fmt.Println(vbf.degree())
}

func test1() {
	fmt.Println("Random (3,2) function")
	vbf, _ := newRandomVBF(3, 2)
	fmt.Println(vbf)
	for i, v := range vbf.getWeight() {
		fmt.Printf("w%d = %d\n", i, v)
	}
	fmt.Println("after Moebius transformation:")
	mvbf := vbf.Moebius()
	fmt.Println(mvbf)
	fmt.Println(mvbf.printANF())
	fmt.Println("\n========================\n")
}

func pow2(n int) int {
	return 1 << n
}
