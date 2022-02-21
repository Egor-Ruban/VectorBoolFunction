package main

import (
	"errors"
	"math/rand"
	"time"
)

type BoolFunction struct {
	value      [][]uint16
	mBlockSize int
	n          int
	m          int
	len        int
}

//вычисляет n из длины
func log2(len int) (i int) {
	for i = 1; (len>>i - 1) != 0; i++ {
	}
	return
}

////подсчитывает вес вектора
//func (b BoolFunction) getWeight() (weight int) {
//	weight = 0
//	for _, v := range b.value {
//		for ; v != 0; v = v & (v - 1) {
//			weight++
//		}
//	}
//	return
//}
//
////интерфейс stringer, для вывода через стандартные функции
//func (b BoolFunction) String() string {
//	res := ""
//	if b.bitSize < 16 {
//		for i := 0; i < b.bitSize; i++ {
//			res += string(((b.value[0] >> (15 - i)) & 1) + '0')
//		}
//	} else {
//		for _, v := range b.value {
//			res += fmt.Sprintf("%016b\n", v)
//		}
//	}
//	return res
//}
//
////красивый вывод в виде таблички
//func (b BoolFunction) printPretty() (weight int) {
//	//todo вывод табличкой
//	return
//}
//
////конструктор по строке
//func newBF(str string) (BoolFunction, error) {
//	len := len(str)
//	if len > 1 && (len&(len-1)) == 0 {
//		var bf BoolFunction
//		if len < 16 {
//			bf = BoolFunction{
//				make([]uint16, 1),
//				1,
//				log2(len),
//				len,
//			}
//		} else {
//			bf = BoolFunction{
//				make([]uint16, len/16),
//				len / 16,
//				log2(len),
//				len,
//			}
//		}
//		for i, v := range str {
//			if v != '0' && v != '1' {
//				return BoolFunction{}, errors.New("wrong vector")
//			}
//			block := i / 16
//			bf.value[block] <<= 1
//			bf.value[block] |= uint16(v - '0')
//		}
//		if len < 16 {
//			bf.value[0] <<= 16 - len
//		}
//		return bf, nil
//	} else {
//		return BoolFunction{}, errors.New("wrong input")
//	}
//}

//конструктор рандомный по длине
func newRandomBF(n, m int) (BoolFunction, error) {
	rand.Seed(time.Now().UnixNano())
	if n > 15 {
		return BoolFunction{}, errors.New("too much variables")
	}
	len := 1 << n

	bf := BoolFunction{
		value:      make([][]uint16, len),
		mBlockSize: (m + 16 - 1) / 16,
		len:        len,
		n:          n,
		m:          m,
	}

	diff := (bf.mBlockSize * 16) - m

	for i := 0; i < len; i++ {
		bf.value[i] = make([]uint16, bf.mBlockSize)
		for j, _ := range bf.value[i] {
			bf.value[i][j] = uint16(rand.Intn(1 << 16))
			if diff > 0 && j == bf.mBlockSize-1 {
				bf.value[i][j] <<= diff
			}
		}
	}
	return bf, nil
}

func newRevBoolFunction(n, m int) (BoolFunction, error) {
	if n < m {

	} else if n == n {

	}
	return BoolFunction{}, errors.New("n > m")
}
