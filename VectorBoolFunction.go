package main

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"
	"unsafe"
)

//not more than 32
type block uint8

type BoolFunction struct {
	value [][]block
	n     int
	m     int
	rows  int

	mBlockSize int
	blockSize  int
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
//	rows := rows(str)
//	if rows > 1 && (rows&(rows-1)) == 0 {
//		var bf BoolFunction
//		if rows < 16 {
//			bf = BoolFunction{
//				make([]uint16, 1),
//				1,
//				log2(rows),
//				rows,
//			}
//		} else {
//			bf = BoolFunction{
//				make([]uint16, rows/16),
//				rows / 16,
//				log2(rows),
//				rows,
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
//		if rows < 16 {
//			bf.value[0] <<= 16 - rows
//		}
//		return bf, nil
//	} else {
//		return BoolFunction{}, errors.New("wrong input")
//	}
//}

func (b block) String() string {
	blockSize := int(unsafe.Sizeof(block(0)) * 8)
	formatString := "%0" + strconv.Itoa(blockSize) + "b"
	return fmt.Sprintf(formatString, b)
}

func (b *block) swap(j int, k int) {
	blockSize := int(unsafe.Sizeof(block(0)) * 8)
	/* Move all bits of first set to rightmost side */
	set1 := (*b >> (blockSize - j - 1)) & 1

	/* Move all bits of second set to rightmost side */
	set2 := (*b >> (blockSize - k - 1)) & 1

	/* Xor the two sets */
	Xor := set1 ^ set2

	/* Put the Xor bits back to their original positions */
	Xor = (Xor << j) | (Xor << k)

	/* Xor the 'Xor' with the original number so that the
	   two sets are swapped */
	result := *b ^ Xor
	*b = result
}

//конструктор рандомный по длине
func newRandomVBF(n, m int) (BoolFunction, error) {
	rand.Seed(time.Now().UnixNano())

	blockSize := int(unsafe.Sizeof(block(0)) * 8)
	if n > blockSize-1 {
		return BoolFunction{}, errors.New("too much variables")
	}
	len := 1 << n

	bf := BoolFunction{
		value:      make([][]block, len),
		mBlockSize: (m + blockSize - 1) / blockSize,
		rows:       len,
		n:          n,
		m:          m,
		blockSize:  blockSize,
	}

	diff := (bf.mBlockSize * blockSize) - m
	fmt.Println("diff = ", diff)
	for i := 0; i < len; i++ {
		bf.value[i] = make([]block, bf.mBlockSize)
		for j, _ := range bf.value[i] {
			bf.value[i][j] = block(rand.Intn(1 << blockSize))
			if diff > 0 && j == bf.mBlockSize-1 {
				bf.value[i][j] <<= diff
			}

			fmt.Print(bf.value[i][j], "|")
		}
		fmt.Println()
	}
	return bf, nil
}

func newRevVBF(n, m int) (BoolFunction, error) {
	rand.Seed(time.Now().UnixNano())
	blockSize := int(unsafe.Sizeof(block(0)) * 8)
	if n > blockSize-1 {
		return BoolFunction{}, errors.New("too much variables")
	}

	len := 1 << n
	bf := BoolFunction{
		value:      make([][]block, len),
		mBlockSize: (m + blockSize - 1) / blockSize,
		rows:       len,
		n:          n,
		m:          m,
		blockSize:  blockSize,
	}
	if n < m {
		//todo
	} else if n == n {
		if bf.mBlockSize > 1 {
			return BoolFunction{}, errors.New("too much variables")
		}
		for i := 0; i < n; i++ {
			bf.value[i] = make([]block, bf.mBlockSize)
			bf.value[i][0] = block(i)
			for j := bf.blockSize; j > 0; j-- {
				k := rand.Intn(j + 1)
				bf.value[i][0].swap(j, k)
			}
		}
		return bf, nil
	}
	return BoolFunction{}, errors.New("n > m")
}
