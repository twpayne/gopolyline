package polyline

import (
	"fmt"
	"strings"
)

type InvalidCharacterError struct {
	pos int
	char byte
}

func (de InvalidCharacterError) Error() string {
	return fmt.Sprintf("invalid character %q at position %d", de.char, de.pos)
}

type UnterminatedError struct {
}

func (ie UnterminatedError) Error() string {
	return "unterminated string"
}

func DecodeUints(s string) ([]uint, error) {
	xs := make([]uint, 0)
	var x, shift uint
	for i, c := range []byte(s) {
		if c < 95 {
			if c < 63 {
				return nil, &InvalidCharacterError{i, c}
			} else {
				xs = append(xs, x+(uint(c)-63)<<shift)
				x = 0
				shift = 0
			}
		} else if c < 127 {
			x += (uint(c) - 95) << shift
			shift += 5
		} else {
			return nil, &InvalidCharacterError{i, c}
		}
	}
	if shift != 0 {
		return nil, &UnterminatedError{}
	}
	return xs, nil
}

func DecodeInts(s string) ([]int, error) {
	xs, err := DecodeUints(s)
	if err != nil {
		return nil, err
	}
	ys := make([]int, len(xs))
	for i, u := range xs {
		if u&1 == 0 {
			ys[i] = int(u >> 1)
		} else {
			ys[i] = -int((u + 1) >> 1)
		}
	}
	return ys, nil
}

func Decode(s string, dim int) ([]float64, error) {
	xs, err := DecodeInts(s)
	if err != nil {
		return nil, err
	}
	ys := make([]float64, len(xs))
	for j, i := range xs {
		ys[j] = float64(i) / 1e5
		if j >= dim {
			ys[j] += ys[j-dim]
		}
	}
	return ys, nil
}

func EncodeUint(x uint) string {
	bs := make([]byte, 0, 7)
	for ; x >= 32; x >>= 5 {
		bs = append(bs, byte((x&31)+95))
	}
	bs = append(bs, byte(x+63))
	return string(bs)
}

func EncodeUints(xs []uint) string {
	ss := make([]string, len(xs))
	for i, x := range xs {
		ss[i] = EncodeUint(x)
	}
	return strings.Join(ss, "")
}

func EncodeInt(x int) string {
	y := uint(x) << 1
	if x < 0 {
		y = ^y
	}
	return EncodeUint(y)
}

func EncodeInts(xs []int) string {
	ss := make([]string, len(xs))
	for i, x := range xs {
		ss[i] = EncodeInt(x)
	}
	return strings.Join(ss, "")
}

func Encode(xs []float64, dim int) string {
	n := len(xs)
	ys := make([]int, n)
	for i := 0; i < n; i++ {
		ys[i] = int(1e5 * xs[i])
	}
	for i := n - 1; i >= dim; i-- {
		ys[i] -= ys[i-dim]
	}
	return EncodeInts(ys)
}
