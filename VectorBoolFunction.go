package main

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"
	"unsafe"
)

type BoolFunction struct {
	value      []blocks
	n          int
	m          int
	rows       int
	mBlockSize int
	blockSize  int

	anf ANF
}

func (bf BoolFunction) getFunction(k int) blocks {
	howManyBlocks := (bf.rows + bf.blockSize - 1) / bf.blockSize
	nf := make(blocks, howManyBlocks)
	for i, v := range bf.value {
		bit := v[k/bf.blockSize] >> (bf.blockSize - k%bf.blockSize - 1) & 1
		nf[i/bf.blockSize] |= bit << (bf.blockSize - i%bf.blockSize - 1)
	}
	return nf
}

//подсчитывает вес вектора
func (bf BoolFunction) getWeight() (weight int) {
	weight = 0
	for _, v := range bf.value {
		for _, v2 := range v {
			for ; v2 != 0; v2 = v2 & (v2 - 1) {
				weight++
			}
		}
	}
	return
}

//интерфейс stringer, для вывода через стандартные функции
func (bf BoolFunction) String() string {
	res := ""
	blockSize := int(unsafe.Sizeof(block(0)) * 8)
	diff := (bf.mBlockSize * bf.blockSize) - bf.m
	formatString := "%0" + strconv.Itoa(blockSize) + "b."         //ordinary format string
	formatString3 := "%0" + strconv.Itoa(bf.blockSize-diff) + "b" //format string if last block is not filled
	for _, v := range bf.value {
		for j, v2 := range v {
			if diff > 0 && j == bf.mBlockSize-1 {
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
func (bf BoolFunction) printPretty() string {
	res := ""
	diff := (bf.mBlockSize * bf.blockSize) - bf.m
	blockSize := int(unsafe.Sizeof(block(0)) * 8)

	formatString := "%0" + strconv.Itoa(blockSize) + "b."         //ordinary format string
	formatString2 := "%0" + strconv.Itoa(bf.n) + "b : "           //format string for variables
	formatString3 := "%0" + strconv.Itoa(bf.blockSize-diff) + "b" //format string if last block is not filled

	for i, v := range bf.value {
		res += fmt.Sprintf(formatString2, i)
		for j, v2 := range v {
			if diff > 0 && j == bf.mBlockSize-1 {
				res += fmt.Sprintf(formatString3, v2>>diff)
			} else {
				res += fmt.Sprintf(formatString, v2)
			}
		}
		res += "\n"
	}
	return res
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
		value:      make([]blocks, len),
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
	bf.anf, _ = generateRandomANF(bf.n, bf.rows, blockSize)
	return bf, nil
}

func newRevVBF(n, m int) (BoolFunction, error) {
	rand.Seed(time.Now().UnixNano())
	blockSize := int(unsafe.Sizeof(block(0)) * 8)
	if n > blockSize-1 {
		return BoolFunction{}, errors.New("too much variables")
	}

	rows := 1 << n
	bf := BoolFunction{
		value:      make([]blocks, rows),
		mBlockSize: (m + blockSize - 1) / blockSize,
		rows:       rows,
		n:          n,
		m:          m,
		blockSize:  blockSize,
	}
	bf.anf, _ = generateRandomANF(bf.n, bf.rows, blockSize)
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
		diff := (bf.mBlockSize * blockSize) - m
		for i := 0; i < bf.rows; i++ {
			bf.value[i] = make([]block, bf.mBlockSize)
			bf.value[i][0] = block(i)
			if diff > 0 {
				bf.value[i][0] <<= diff
			}
		}
		for i := bf.rows - 1; i > 0; i-- {
			j := rand.Intn(i + 1)
			t := bf.value[i][0]
			bf.value[i][0] = bf.value[j][0]
			bf.value[j][0] = t
		}
		return bf, nil
	}
	return BoolFunction{}, errors.New("n > m")
}
