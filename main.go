package main

import "fmt"

func main() {
	bf, _ := newRandomVBF(3, 4)
	fmt.Println(bf.printPretty())
	fmt.Println("weight is ", bf.getWeight())
}
