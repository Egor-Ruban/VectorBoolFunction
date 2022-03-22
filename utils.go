package main

//вычисляет n из длины
func log2(len int) (i int) {
	for i = 1; (len>>i - 1) != 0; i++ {
	}
	return
}

func Combinations(n, k int) int {
	nf := 1
	kf := 1
	nkf := 1

	for i := 1; i <= nf; i++ {
		nf *= i
	}
	for i := 1; i <= kf; i++ {
		kf *= i
	}
	for i := 1; i <= nkf; i++ {
		nkf *= i
	}
	return nf / (kf * nkf)
}
