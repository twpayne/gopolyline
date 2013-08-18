PACKAGE DOCUMENTATION

package polyline
    import "github.com/twpayne/polyline"

    Package polyline provides functions for Google Maps' Encoded Polyline
    Algorithm Format. See
    https://developers.google.com/maps/documentation/utilities/polylinealgorithm


FUNCTIONS

func Decode(s string, dim int) ([]float64, error)
    Decode decodes a polyline of dimension dim from a string. The polyline
    is returned as a flat slice, e.g.

	[]float64{lat1, lng1, lat2, lng2...}

    The dimension dim should normally be 2.

func DecodeInts(s string) ([]int, error)
    DecodeInts decodes a slice of ints from a string.

func DecodeUints(s string) ([]uint, error)
    DecodeUnits decodes a slice of uints from a string.

func Encode(xs []float64, dim int) string
    Encode encodes a polyling as a string. The polyline must be structured
    as a flat slice, e.g.

	[]float64{lat1, lng1, lat2, lng2...}

    The dimension dim should normally be 2.

func EncodeInt(x int) string
    EncodeInt encodes a single int as a string.

func EncodeInts(xs []int) string
    EncodeInts encodes a slice of ints as a string.

func EncodeUint(x uint) string
    EncodeUnit encodes a single uint as a string.

func EncodeUints(xs []uint) string
    EncodeUnits encodes a slice of uints as a string.


TYPES

type InvalidCharacterError struct {
    // contains filtered or unexported fields
}
    InvalidCharacterError is returned when an invalid character is
    encountered.


func (de InvalidCharacterError) Error() string


type UnterminatedError struct {
}
    UnterminatedError is returned when the string is unterminated.


func (ie UnterminatedError) Error() string



