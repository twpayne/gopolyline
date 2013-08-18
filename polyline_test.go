package polyline

import (
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
)

func TestDecodeInts(t *testing.T) {
	tests := []struct {
		s    string
		want []int
	}{
		{"_p~iF", []int{3850000}},
		{"~ps|U", []int{-12020000}},
		{"_p~iF~ps|U", []int{3850000, -12020000}},
		{"_p~iF~ps|U_ulLnnqC_mqNvxq`@", []int{3850000, -12020000, 220000, -75000, 255200, -550300}},
	}
	for _, c := range tests {
		if got, err := DecodeInts(c.s); !reflect.DeepEqual(got, c.want) || err != nil {
			t.Errorf("DecodeInts(%q) == %q, %q, want %q, nil", c.s, got, err, c.want)
		}
	}
}

func TestDecodeIntsErrors(t *testing.T) {
	tests := []struct {
		s    string
		want string
	}{
		{"0", "invalid character '0' at position 0"},
		{"_p~iF~ps|U ", "invalid character ' ' at position 10"},
		{"_p~i", "unterminated string"},
		{"_p~iF~ps|u", "unterminated string"},
	}
	for _, c := range tests {
		if _, err := DecodeInts(c.s); err == nil || err.Error() != c.want {
			t.Errorf("DecodeInts(%q) error = %q, want %q", c.s, err, c.want)
		}
	}
}

func TestDecode(t *testing.T) {
	tests := []struct {
		s      string
		stride int
		want   []float64
	}{
		{"_p~iF~ps|U", 2, []float64{38.5, -120.2}},
		{"_p~iF~ps|U_ulLnnqC_mqNvxq`@", 2, []float64{38.5, -120.2, 40.7, -120.95, 43.252, -126.453}},
	}
	for _, c := range tests {
		if got, err := Decode(c.s, c.stride); !reflect.DeepEqual(got, c.want) || err != nil {
			t.Errorf("Decode(%q) == %q, %q, want %q, nil", c.s, got, err, c.want)
		}
	}
}

func TestEncodeUint(t *testing.T) {
	tests := []struct {
		x    uint
		want string
	}{
		{0, "?"},
		{1, "@"},
		{10, "I"},
		{100, "cB"},
		{174, "mD"},
		{1000, "g^"},
		{10000, "owH"},
		{100000, "_t`B"},
		{1000000, "_qo]"},
	}
	for _, c := range tests {
		if got := EncodeUint(c.x); got != c.want {
			t.Errorf("EncodeUint(%q) == %q, want %q", c.x, got, c.want)
		}
	}
}

func TestEncodeInt(t *testing.T) {
	tests := []struct {
		x    int
		want string
	}{
		{-1000000, "~b`|@"},
		{-100000, "~hbE"},
		{-10000, "~oR"},
		{-1000, "n}@"},
		{-100, "fE"},
		{-10, "R"},
		{-1, "@"},
		{0, "?"},
		{1, "A"},
		{10, "S"},
		{100, "gE"},
		{1000, "o}@"},
		{10000, "_pR"},
		{100000, "_ibE"},
		{1000000, "_c`|@"},
	}
	for _, c := range tests {
		if got := EncodeInt(c.x); got != c.want {
			t.Errorf("EncodeInt(%q) == %q, want %q", c.x, got, c.want)
		}
	}
}

func TestEncode(t *testing.T) {
	tests := []struct {
		xs     []float64
		stride int
		want   string
	}{
		{[]float64{38.5, -120.2}, 2, "_p~iF~ps|U"},
		{[]float64{38.5, -120.2, 40.7, -120.95, 43.252, -126.453}, 2, "_p~iF~ps|U_ulLnnqC_mqNvxq`@"},
	}
	for _, c := range tests {
		if got := Encode(c.xs, c.stride); got != c.want {
			t.Errorf("Encode(%g, %d) == %q, want %q", c.xs, c.stride, got, c.want)
		}
	}
}

func TestIntsEncodeDecodeQuick(t *testing.T) {
	f := func(xs []int) bool {
		ys, err := DecodeInts(EncodeInts(xs))
		return reflect.DeepEqual(xs, ys) && err == nil
	}
	quick.Check(f, nil)
}

func TestUintsEncodeDecodeQuick(t *testing.T) {
	f := func(xs []uint) bool {
		ys, err := DecodeUints(EncodeUints(xs))
		return reflect.DeepEqual(xs, ys) && err == nil
	}
	quick.Check(f, nil)
}

const complexSize = 50

type encodedString string

func (encodedString) Generate(rand *rand.Rand, size int) reflect.Value {
	numChars := rand.Intn(complexSize)
	bytes := make([]byte, numChars)
	for i := 0; i < numChars; i++ {
		bytes[i] = byte(63 + rand.Intn(64))
	}
	if numChars > 0 {
		bytes[numChars-1] |= 0x20
	}
	return reflect.ValueOf(encodedString(bytes))
}

func TestDecodeEncodeUintsQuick(t *testing.T) {
	f := func(es encodedString) bool {
		s := string(es)
		xs, err := DecodeUints(s)
		return err == nil && EncodeUints(xs) == s
	}
	quick.Check(f, nil)
}

func TestDecodeEncodeIntsQuick(t *testing.T) {
	f := func(es encodedString) bool {
		s := string(es)
		xs, err := DecodeInts(s)
		return err == nil && EncodeInts(xs) == s
	}
	quick.Check(f, nil)
}

func TestDecodeEncodeQuick(t *testing.T) {
	f := func(es encodedString) bool {
		s := string(es)
		xs, err := Decode(s, 2)
		return err == nil && Encode(xs, 2) == s
	}
	quick.Check(f, nil)
}
