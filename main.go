package main

import "fmt"

func main() {
	bf, _ := newRandomVBF(3, 4)
	fmt.Println(bf.printPretty())
	//fmt.Println("weight is ", bf.getWeight())
	fmt.Println("anf vector:", bf.anf.printVector())
	fmt.Println("anf equation:", bf.anf)
	//s := "123\n"
	//fmt.Println(string('z' - 1))
}
