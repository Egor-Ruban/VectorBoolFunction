package main

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type ANF struct {
	n              int
	rows           int
	blockSize      int
	vector         blocks
	blockLengthANF int
}

func (anf ANF) String() string {
	if anf.n > anf.blockSize {
		return ""
	}
	res := ""
	if (anf.vector[0]>>(anf.blockSize-1))&1 == 1 {
		res += "1"
	}
	for i := 1; i < anf.rows; i++ {
		if (anf.vector[i/anf.blockSize]>>(anf.blockSize-i%anf.blockSize-1))&1 == 1 {
			if len(res) != 0 {
				res += "+"
			}
			for j := 0; j < anf.n; j++ {
				if (i>>j)&1 == 1 {
					res += string(rune('z' - j))
				}
			}
		}
	}
	if len(res) == 0 {
		return "0"
	}
	return res
}

func (anf ANF) printVector() string {
	res := ""
	diff := (anf.blockLengthANF * anf.blockSize) - anf.rows
	formatString := "%0" + strconv.Itoa(anf.blockSize) + "b."      //ordinary format string
	formatString3 := "%0" + strconv.Itoa(anf.blockSize-diff) + "b" //format string if last block is not filled

	for j, v2 := range anf.vector {
		if diff > 0 && j == anf.blockLengthANF-1 {
			res += fmt.Sprintf(formatString3, v2>>diff)
		} else {
			res += fmt.Sprintf(formatString, v2)
		}
	}

	return res
}

func generateRandomANF(n, rows, blockSize int) (ANF, error) {
	rand.Seed(time.Now().UnixNano())
	if n > blockSize {
		return ANF{}, errors.New("too many variables")
	}
	anf := ANF{
		n:              n,
		rows:           rows,
		blockSize:      blockSize,
		blockLengthANF: (rows + blockSize - 1) / blockSize,
		vector:         make([]block, (rows+blockSize-1)/blockSize),
	}
	diff := (anf.blockLengthANF * anf.blockSize) - anf.rows
	for j, _ := range anf.vector {
		anf.vector[j] = block(rand.Intn(1 << anf.blockSize))
		if diff > 0 && j == anf.blockLengthANF-1 {
			anf.vector[j] <<= diff
		}
	}

	return anf, nil
}

func Moebius(function []block, rows int, blockSize int) ANF {
	anf := ANF{
		n:              log2(rows),
		rows:           rows,
		blockSize:      blockSize,
		blockLengthANF: (rows + blockSize - 1) / blockSize,
		vector:         make([]block, (rows+blockSize-1)/blockSize),
	}

	for i := range function {
		anf.vector[i] = function[i]
	}
	n := log2(rows)
	for i := 0; i < n; i++ {
		anf.vector = anf.vector.xor(anf.vector.shiftRight(1 << i).and(generateMask(anf.blockLengthANF, i)))
	}
	return anf
}
