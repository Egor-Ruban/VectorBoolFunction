package main

import "fmt"

func main() {
	//for m := 1; m < 5; m++ {
	//	for n := 2; n < 20; n++ {
	//		bf, _ := newRandomVBF(n, m)
	//		//fmt.Println("\n", bf.printPretty())
	//		w := bf.getWeight()
	//		fmt.Println("n =", n, "m =", m, " weight =", w, " k =", float64(w)/(float64(1<<n)*float64(m)))
	//	}
	//}
	bf, _ := newRevVBF(3, 4)
	fmt.Println(bf.printPretty())
	//fmt.Println("weight is ", bf.getWeight())
	//fmt.Println("anf vector:", bf.anf.printVector())
	//fmt.Println("anf equation:", bf.anf)
	//s := "123\n"
	//fmt.Println(string('z' - 1))
}
