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
			moebius2 := moebius.Moebius()
			if n == 2 {
				fmt.Println("function:")
				fmt.Println(l)
				fmt.Println("one moebius: ")
				fmt.Println(moebius)
				fmt.Println("two moebius: ")
				fmt.Println(moebius2)
			}
			fmt.Println("moebius(moebius(vbf)) == vbf :", l.isEqual(moebius2))
			fmt.Println("=============================================")
		}
	}
}
