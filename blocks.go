package main

import (
	"fmt"
	"strconv"
	"unsafe"
)

type blocks []block

func generateMask(howManyBlocks, step int) blocks {
	blockSize := int(unsafe.Sizeof(block(0)) * 8)
	mask := make(blocks, howManyBlocks)
	if (1 << step) < blockSize {
		m := block(0)
		bit := block(0)
		for pushed := 0; pushed < blockSize; {
			for i := 0; i < (1 << step); i++ {
				m <<= 1
				m |= bit
				pushed++
			}
			bit ^= 1
		}
		for i := range mask {
			mask[i] = m
		}
	} else {
		m := block(1<<blockSize - 1)
		for i := range mask {
			if (i*blockSize)%(1<<step) == 0 {
				m ^= block(1<<blockSize - 1)
			}
			mask[i] = m
		}
	}
	return mask
}

func (b blocks) shiftLeft(n int) blocks {
	blockSize := int(unsafe.Sizeof(block(0)) * 8)
	t := make(blocks, len(b))
	if n < blockSize {
		t[0] = b[0] << n
		for i := 1; i < len(b); i++ {
			t[i] = b[i]
			saved := t[i] >> (blockSize - n)
			t[i-1] |= saved
			t[i] <<= n
		}
	} else if n%blockSize == 0 {
		k := n / blockSize
		for i := k; i < len(b); i++ {
			t[i-k] = b[i]
		}
	} else if n < len(b)*blockSize {
		k := n / blockSize
		for i := k; i < len(b); i++ {
			t[i-k] = b[i]
		}
		k = n % blockSize
		t[0] = t[0] << k
		for i := 1; i < len(b); i++ {
			saved := t[i] >> (blockSize - k)
			t[i-1] |= saved
			t[i] <<= k
		}
	}
	return t
}

func (b blocks) shiftRight(n int) blocks {
	blockSize := int(unsafe.Sizeof(block(0)) * 8)
	t := make(blocks, len(b))
	if n < blockSize {
		t[len(b)-1] = b[len(b)-1] >> n
		for i := len(b) - 2; i >= 0; i-- {
			t[i] = b[i]
			saved := t[i] << (blockSize - n)
			t[i+1] |= saved
			t[i] >>= n
		}
	} else if n%blockSize == 0 {
		k := n / blockSize
		for i := len(b) - 1 - k; i >= 0; i-- {
			t[i+k] = b[i]
		}
	} else if n < len(b)*blockSize {
		k := n / blockSize
		for i := len(b) - 1 - k; i >= 0; i-- {
			t[i+k] = b[i]
		}
		n = n % blockSize
		t[len(b)-1] = t[len(b)-1] >> n
		for i := len(b) - 2; i >= 0; i-- {
			saved := t[i] << (blockSize - n)
			t[i+1] |= saved
			t[i] >>= n
		}
	}
	return t
}

func (b blocks) String() string {
	res := ""
	blockSize := int(unsafe.Sizeof(block(0)) * 8)
	for _, v := range b {
		formatString := "%0" + strconv.Itoa(blockSize) + "b."
		res += fmt.Sprintf(formatString, v)
	}
	return res
}

func (b blocks) or(b2 blocks) blocks {
	t := make(blocks, len(b))
	for i := range b {
		t[i] = b[i] | b2[i]
	}
	return t
}

func (b blocks) xor(b2 blocks) blocks {
	t := make(blocks, len(b))
	for i := range b {
		t[i] = b[i] ^ b2[i]
	}
	return t
}

func (b blocks) and(b2 blocks) blocks {
	t := make(blocks, len(b))
	for i := range b {
		t[i] = b[i] & b2[i]
	}
	return t
}
