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
type block uint16

type BoolFunction struct {
	value [][]block
	n     int
	m     int
	rows  int

	mBlockSize int
	blockSize  int
}

//подсчитывает вес вектора
func (b BoolFunction) getWeight() (weight int) {
	weight = 0
	for _, v := range b.value {
		for _, v2 := range v {
			for ; v2 != 0; v2 = v2 & (v2 - 1) {
				weight++
			}
		}
	}
	return
}

//интерфейс stringer, для вывода через стандартные функции
func (b BoolFunction) String() string {
	res := ""
	blockSize := int(unsafe.Sizeof(block(0)) * 8)
	diff := (b.mBlockSize * b.blockSize) - b.m
	formatString := "%0" + strconv.Itoa(blockSize) + "b."        //ordinary format string
	formatString3 := "%0" + strconv.Itoa(b.blockSize-diff) + "b" //format string if last block is not filled
	for _, v := range b.value {
		for j, v2 := range v {
			if diff > 0 && j == b.mBlockSize-1 {
				res += fmt.Sprintf(formatString3, v2>>diff)
			} else {
				res += fmt.Sprintf(formatString, v2)
			}
		}
		res += "\n"
	}
	return res
}

//красивый вывод в виде таблички
func (b BoolFunction) printPretty() string {
	res := ""
	diff := (b.mBlockSize * b.blockSize) - b.m
	blockSize := int(unsafe.Sizeof(block(0)) * 8)

	formatString := "%0" + strconv.Itoa(blockSize) + "b."        //ordinary format string
	formatString2 := "%0" + strconv.Itoa(b.n) + "b : "           //format string for variables
	formatString3 := "%0" + strconv.Itoa(b.blockSize-diff) + "b" //format string if last block is not filled

	for i, v := range b.value {
		res += fmt.Sprintf(formatString2, i)
		for j, v2 := range v {
			if diff > 0 && j == b.mBlockSize-1 {
				res += fmt.Sprintf(formatString3, v2>>diff)
			} else {
				res += fmt.Sprintf(formatString, v2)
			}
		}
		res += "\n"
	}
	return res
}

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
	set1 := (*b >> (blockSize - j - 1)) & 1
	set2 := (*b >> (blockSize - k - 1)) & 1
	Xor := set1 ^ set2
	Xor = (Xor << j) | (Xor << k)
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
	for i := 0; i < len; i++ {
		bf.value[i] = make([]block, bf.mBlockSize)
		for j, _ := range bf.value[i] {
			bf.value[i][j] = block(rand.Intn(1 << blockSize))
			if diff > 0 && j == bf.mBlockSize-1 {
				bf.value[i][j] <<= diff
			}
		}
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
		if m > blockSize {
			return BoolFunction{}, errors.New("m is too big")
		}
		l := ((uint64(1) << m) + 64 - 1) / 64 //[2^m/64] - сколько блоков по 64 надо для 2^m бит
		isTaken := make([]uint64, l)
		for i := range isTaken {
			isTaken[i] = 0
		}

		for i := 0; i < bf.rows; i++ {
			x := rand.Intn(1 << m)
			if (isTaken[x/64]>>(64-x%64-1))&1 == 1 {
				i--
			} else {
				bf.value[i] = make([]block, 1)
				bf.value[i][0] = block(x << (blockSize - m))
				isTaken[x/64] |= uint64(1) << (64 - x%64 - 1)
			}
		}
		return bf, nil
	} else if n == n {
		for i := 0; i < bf.rows; i++ {
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
