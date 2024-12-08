package jsonx

import (
	"strconv"
)

// Number represents a JSON number
type Number struct {
	hasDecimals bool
	str         string
}

// String returns the Go string version of it, maybe you want arbitrary size numbers
func (n Number) String() string {
	return n.str
}

// HasDecimals - Returns if the number has decimal digits
func (n Number) HasDecimals() bool {
	return n.hasDecimals
}

// AsFloat64 - Returns it as a float64
func (n Number) AsFloat64() (float64, error) {
	return strconv.ParseFloat(n.str, 64)
}
