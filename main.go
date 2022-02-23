package main

import "fmt"

func main() {
	for n := 2; n < 7; n++ {
		fmt.Println("n   = ", n)
		bf, _ := newRandomVBF(n, 1)
		t := bf.getFunction(0)
		anf := Moebius(t, bf.rows, bf.blockSize)
		t2 := Moebius(anf.vector, anf.rows, anf.blockSize)
		//fmt.Println(bf.printPretty())

		fmt.Println("t   = ", t)
		fmt.Println("anf = ", anf.vector)
		fmt.Println("t2  = ", t2.vector)
		fmt.Println("===================")
	}

}
