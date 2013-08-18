// Package polyline provides functions for Google Maps' Encoded Polyline Algorithm Format.
// See https://developers.google.com/maps/documentation/utilities/polylinealgorithm
package polyline

import (
	"fmt"
	"strings"
)

// InvalidCharacterError is returned when an invalid character is encountered.
type InvalidCharacterError struct {
	pos int
	char byte
}

func (de InvalidCharacterError) Error() string {
	return fmt.Sprintf("invalid character %q at position %d", de.char, de.pos)
}

// UnterminatedError is returned when the string is unterminated.
type UnterminatedError struct {
}

func (ie UnterminatedError) Error() string {
	return "unterminated string"
}

// DecodeUnits decodes a slice of uints from a string.
func DecodeUints(s string) ([]uint, error) {
	xs := make([]uint, 0)
	var x, shift uint
	for i, c := range []byte(s) {
		switch {
		case 63 <= c && c < 95:
			xs = append(xs, x+(uint(c)-63)<<shift)
			x = 0
			shift = 0
		case 95 <= c && c < 127:
			x += (uint(c) - 95) << shift
			shift += 5
		default:
			return nil, &InvalidCharacterError{i, c}
		}
	}
	if shift != 0 {
		return nil, &UnterminatedError{}
	}
	return xs, nil
}

// DecodeInts decodes a slice of ints from a string.
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

// Decode decodes a polyline of dimension dim from a string.
// The polyline is returned as a flat slice, e.g.
//     []float64{lat1, lng1, lat2, lng2...}
// The dimension dim should normally be 2.
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

// EncodeUnit encodes a single uint as a string.
func EncodeUint(x uint) string {
	bs := make([]byte, 0, 7)
	for ; x >= 32; x >>= 5 {
		bs = append(bs, byte((x&31)+95))
	}
	bs = append(bs, byte(x+63))
	return string(bs)
}

// EncodeUnits encodes a slice of uints as a string.
func EncodeUints(xs []uint) string {
	ss := make([]string, len(xs))
	for i, x := range xs {
		ss[i] = EncodeUint(x)
	}
	return strings.Join(ss, "")
}

// EncodeInt encodes a single int as a string.
func EncodeInt(x int) string {
	y := uint(x) << 1
	if x < 0 {
		y = ^y
	}
	return EncodeUint(y)
}

// EncodeInts encodes a slice of ints as a string.
func EncodeInts(xs []int) string {
	ss := make([]string, len(xs))
	for i, x := range xs {
		ss[i] = EncodeInt(x)
	}
	return strings.Join(ss, "")
}

// Encode encodes a polyling as a string.
// The polyline must be structured as a flat slice, e.g.
//     []float64{lat1, lng1, lat2, lng2...}
// The dimension dim should normally be 2.
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
