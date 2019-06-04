package modesdecoder

import (
	"fmt"
	"strconv"
	"strings"
)

type Bits struct {
	bits []uint8
	err  error
}

func NewBits(bits []uint8) *Bits {
	return &Bits{bits: bits}
}

func ParseHex(hex string) (*Bits, error) {
	var bin Bits
	for _, r := range strings.ToLower(hex) {
		c, ok := hexToBinLookup[r]
		if !ok {
			return nil, fmt.Errorf("invalid hex: %v", c)
		}
		bin.bits = append(bin.bits, c...)
	}
	return &bin, nil
}

func (b *Bits) Len() int {
	return len(b.bits)
}

func (b *Bits) Err() error {
	return b.err
}

func (b *Bits) Copy() *Bits {
	bits := make([]uint8, len(b.bits))
	copy(bits, b.bits)
	return &Bits{bits: bits}
}

func (b *Bits) Raw() []uint8 { return b.bits }

func (b *Bits) At(i int) uint8 {
	if b.err != nil {
		return 0
	}
	if i < 0 || i >= len(b.bits) {
		b.err = fmt.Errorf("bits out of range (i=%d, len=%d) ", i, len(b.bits))
		return 0
	}
	return b.bits[i]
}

func (b *Bits) Bool(i int) bool {
	return b.At(i) == 1
}

func (b *Bits) String(from, to int) string {
	if b.err != nil {
		return ""
	}
	s := ""
	for _, x := range b.slice(from, to) {
		s += strconv.Itoa(int(x))
	}
	return s
}

func (b *Bits) Int64(from, to int) int64 {
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

func (b *Bits) Uint(from, to int) uint {
	return uint(b.Int64(from, to))
}

func (b *Bits) Slice(from, to int) *Bits {
	bits := b.slice(from, to)
	return &Bits{bits: bits, err: b.err}
}

func (b *Bits) slice(from, to int) []uint8 {
	if b.err != nil {
		return nil
	}
	if from < 0 || to >= len(b.bits) {
		b.err = fmt.Errorf("bits out of range (from=%d, to=%d, len=%d) ", from, to, len(b.bits))
		return nil
	}
	return b.bits[from:to]
}
