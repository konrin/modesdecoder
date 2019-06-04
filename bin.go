package modesdecoder

import (
	"fmt"
	"strconv"
	"strings"
)

type Bin struct {
	bits []uint8
	err  error
}

func (b *Bin) Len() int {
	return len(b.bits)
}

func (b *Bin) Err() error {
	return b.err
}

func ParseHex(hex string) (*Bin, error) {
	var bin Bin
	for _, r := range strings.ToLower(hex) {
		c, ok := hexToBinLookup[r]
		if !ok {
			return nil, fmt.Errorf("invalid hex: %v", c)
		}
		bin.bits = append(bin.bits, c...)
	}
	return &bin, nil
}

func (b *Bin) Bits(from, to int) []uint8 {
	if b.err != nil {
		return nil
	}
	if from < 0 || to >= len(b.bits) {
		b.err = fmt.Errorf("bits out of range (from=%d, to=%d, len=%d) ", from, to, len(b.bits))
		return nil
	}
	return b.bits[from:to]
}

func (b *Bin) String(from, to int) string {
	if b.err != nil {
		return ""
	}
	s := ""
	for _, x := range b.Bits(from, to) {
		s += strconv.Itoa(int(x))
	}
	return s
}

func (b *Bin) Int64(from, to int) int64 {
	s := b.String(from, to)
	if b.err != nil {
		return 0
	}
	x, err := strconv.ParseInt(s, 2, 64)
	if err != nil {
		b.err = err
		return 0
	}
	return x
}
