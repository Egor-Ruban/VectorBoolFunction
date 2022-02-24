package main

import "fmt"

func main() {
	l, _ := newRevVBF(2, 4)
	fmt.Println(l.Moebius())
	fmt.Println(l.Moebius().printANF())
}
