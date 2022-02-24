package main

import (
	"fmt"
	"strconv"
	"unsafe"
)

//not more than 32
type block uint16

func blockSize() int {
	return int(unsafe.Sizeof(block(0)) * 8)
}

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
