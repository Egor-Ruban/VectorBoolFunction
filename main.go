package main

import (
	"fmt"
)

func main() {
	for n := 2; n < 20; n++ {
		for m := 1; m < 4; m++ {
			l, _ := newRandomVBF(n, m)
			fmt.Println("n =", n, ", m =", m)
			moebius := l.Moebius()
			newmoebius := l.newMoebius()
			fmt.Println("moebius(vbf) == newMoebius(vbf) :", moebius.isEqual(newmoebius))
			fmt.Println("=============================================")
		}
	}
}
