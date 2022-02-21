package main

import "fmt"

func main() {
	bf, _ := newRandomVBF(4, 20)
	fmt.Println(bf.printPretty())
}
