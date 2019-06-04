package modesdecoder

import (
	"fmt"
	"strconv"
)

type Bin struct {
	bits []uint8
	err  error
}

func (b *Bin) Err() error {
	return b.err
}

func ParseBin(s string) (*Bin, error) {
	var bin Bin
	for _, ch := range s {
		i, err := strconv.Atoi(string(ch))
		if err != nil {
			return nil, err
		}
		bin.bits = append(bin.bits, uint8(i))
	}
	return &bin, nil
}

func (b *Bin) Bits(from, to int) []uint8 {
	if b.err != nil {
		return nil
	}
	if from <= 0 || to >= len(b.bits) {
		b.err = fmt.Errorf("bits out of range (%d - %d)", from, to)
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
