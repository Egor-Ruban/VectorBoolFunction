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
type block uint32

func blockSize() int {
	return int(unsafe.Sizeof(block(0)) * 8)
}

type BoolFunction struct {
	value     []block
	n         int
	m         int
	rows      int
	wasteBits int
}

//подсчитывает вес вектора
func (bf BoolFunction) getWeight() []int {
	weight := make([]int, bf.m)
	for _, v := range bf.value {
		for i := 0; i < bf.m; i++ {
			weight[i] += int(v>>(blockSize()-1-i)) & 1
		}
	}
	return weight
}

//красивый вывод в виде таблички
func (bf BoolFunction) String() string {
	res := ""
	diff := blockSize() - bf.m
	formatString2 := "%0" + strconv.Itoa(bf.n) + "b : "            //format string for variables
	formatString3 := "%0" + strconv.Itoa(blockSize()-diff) + "b\n" //format string if last block is not filled

	for i, v := range bf.value {
		res += fmt.Sprintf(formatString2, i) //можно убрать, если не нужен вывод значений переменных
		res += fmt.Sprintf(formatString3, v>>diff)
	}
	return res
}

//конструктор рандомный по длине
func newRandomVBF(n, m int) (BoolFunction, error) {
	rand.Seed(time.Now().UnixNano())

	if n > blockSize() || m > blockSize() {
		return BoolFunction{}, errors.New("n or m is too big")
	}

	bf := BoolFunction{
		value:     make([]block, 1<<n),
		rows:      1 << n,
		n:         n,
		m:         m,
		wasteBits: blockSize() - m,
	}

	for i := 0; i < bf.rows; i++ {
		bf.value[i] = block(rand.Intn(1<<m)) << bf.wasteBits
	}
	return bf, nil
}

func newRevVBF(n, m int) (BoolFunction, error) {
	rand.Seed(time.Now().UnixNano())
	if n > blockSize() || m > blockSize() {
		return BoolFunction{}, errors.New("n or m is too big")
	}

	bf := BoolFunction{
		value:     make([]block, 1<<n),
		rows:      1 << n,
		n:         n,
		m:         m,
		wasteBits: blockSize() - m,
	}

	if n < m {
		l := ((uint64(1) << m) + 64 - 1) / 64 //[2^m/64] - сколько блоков по 64 надо для 2^m бит
		isTaken := make([]uint64, l)

		for i := 0; i < bf.rows; i++ {
			x := rand.Intn(1 << m)
			if (isTaken[x/64]>>(64-x%64-1))&1 == 1 {
				i--
			} else {
				bf.value[i] = block(x << bf.wasteBits)
				isTaken[x/64] |= uint64(1) << (64 - x%64 - 1)
			}
		}
		return bf, nil
	} else if n == n {
		for i := 0; i < bf.rows; i++ {
			bf.value[i] = block(i) << bf.wasteBits
		}
		for i := bf.rows - 1; i > 0; i-- {
			j := rand.Intn(i + 1)
			t := bf.value[i]
			bf.value[i] = bf.value[j]
			bf.value[j] = t
		}
		return bf, nil
	}
	return BoolFunction{}, errors.New("n > m")
}

func (bf BoolFunction) shiftDown(k int) BoolFunction {
	t := BoolFunction{
		value:     make([]block, bf.rows),
		rows:      bf.rows,
		n:         bf.n,
		m:         bf.m,
		wasteBits: blockSize() - bf.m,
	}

	for i := bf.rows - 1; i >= k; i-- {
		t.value[i] = bf.value[i-k]
	}
	return t
}

func (bf BoolFunction) shiftUp(k int) BoolFunction {
	t := BoolFunction{
		value:     make([]block, bf.rows),
		rows:      bf.rows,
		n:         bf.n,
		m:         bf.m,
		wasteBits: blockSize() - bf.m,
	}

	for i := 0; i < bf.rows-k; i++ {
		t.value[i] = bf.value[i+k]
	}
	return t
}

func (bf BoolFunction) xor(bf2 BoolFunction) BoolFunction {
	t := BoolFunction{
		value:     make([]block, bf.rows),
		rows:      bf.rows,
		n:         bf.n,
		m:         bf.m,
		wasteBits: blockSize() - bf.m,
	}

	for i := 0; i < bf.rows; i++ {
		t.value[i] = bf.value[i] ^ bf2.value[i]
	}
	return t
}

func (bf BoolFunction) and(bf2 BoolFunction) BoolFunction {
	t := BoolFunction{
		value:     make([]block, bf.rows),
		rows:      bf.rows,
		n:         bf.n,
		m:         bf.m,
		wasteBits: blockSize() - bf.m,
	}

	for i := 0; i < bf.rows; i++ {
		t.value[i] = bf.value[i] & bf2.value[i]
	}
	return t
}

func (bf BoolFunction) or(bf2 BoolFunction) BoolFunction {
	t := BoolFunction{
		value:     make([]block, bf.rows),
		rows:      bf.rows,
		n:         bf.n,
		m:         bf.m,
		wasteBits: blockSize() - bf.m,
	}

	for i := 0; i < bf.rows; i++ {
		t.value[i] = bf.value[i] | bf2.value[i]
	}
	return t
}

func (bf BoolFunction) generateMask(step int) BoolFunction {
	mask := BoolFunction{
		value:     make([]block, bf.rows),
		rows:      bf.rows,
		n:         bf.n,
		m:         bf.m,
		wasteBits: blockSize() - bf.m,
	}
	m := block(1<<blockSize() - 1)
	for i := range mask.value {
		if (i % (1 << step)) == 0 {
			m ^= block(1<<blockSize() - 1)
		}
		mask.value[i] = m
	}
	return mask
}

func (bf BoolFunction) Moebius() BoolFunction {
	anf := BoolFunction{
		value:     make([]block, bf.rows),
		rows:      bf.rows,
		n:         bf.n,
		m:         bf.m,
		wasteBits: blockSize() - bf.m,
	}

	for i := range bf.value {
		anf.value[i] = bf.value[i]
	}
	for i := 0; i < anf.n; i++ {
		anf = anf.xor(anf.shiftDown(1 << i).and(anf.generateMask(i)))
	}
	return anf
}

func (bf BoolFunction) printANF() string {
	if bf.n > 26 {
		return "too many variables in ANF"
	}

	ANFs := make([]string, bf.m)
	for j := 0; j < bf.m; j++ {
		if (bf.value[0]>>(blockSize()-j-1))&1 == 1 {
			ANFs[j] += "1"
		}
	}
	for i := 1; i < bf.rows; i++ {
		for j := 0; j < bf.m; j++ {
			if (bf.value[i]>>(blockSize()-j-1))&1 == 1 {
				if len(ANFs[j]) != 0 {
					ANFs[j] += "+"
				}
				for k := 0; k < bf.n; k++ {
					if (i>>k)&1 == 1 {
						ANFs[j] += string(rune('z' - k))
					}
				}
			}
		}
	}

	res := ""
	for i := range ANFs {
		res += "anf" + strconv.Itoa(i) + " = "
		if len(ANFs[i]) == 0 {
			res += "0"
		} else {
			res += ANFs[i]
		}
		res += "\n"
	}

	return res
}
