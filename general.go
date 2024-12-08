package jsonx

import (
	"errors"
)

// ErrDefault - Default error for Decode()
var ErrDefault = errors.New("JSON could not be parsed")

// ErrInvalidNumber - Number wrongly formatted
var ErrInvalidNumber = errors.New("invalid number format")

// state - Internal structure for keeping track of state
type state struct {
	src []byte // The whole input
	pos int    // Current position in source

	returnFloats bool
}

// Object represents a JSON Object
type Object = map[string]any

// Array - Represents a JSON Array
type Array = []any

// Identifiers
const (
	// EOS - Used internally to signify end of stream
	iEOS byte = 0x03 // 0x03 = End of Text, 0x04 = End of Transmission
	// Star - Used internally to support jsonc: JSON Comments
	iStar byte = '*'
	// Slash - Used internally to support jsonc: JSON Comments
	iFSlash byte = '/'

	// Dot - Used internally for reading floats
	iDot byte = '.'
	// Quotation - Used internally for reading strings
	iQuotation byte = '"'

	// Hyphen - Syntax literal negative
	iHyphen byte = '-'

	// Comma - Syntax literal comma
	iComma byte = ','
	// Colon - Syntax literal colon
	iColon byte = ':'

	// LeftBrace - Syntax literal to start an object
	iLeftBrace byte = '{'
	// RightBrace - Syntax literal to end an object
	iRightBrace byte = '}'

	// LeftBracket - Syntax literal to start a list
	iLeftBracket byte = '['
	// RightBracket - Syntax literal to end a list
	iRightBracket byte = ']'
)

// get returns the byte at a given index in source
func (state *state) get(n int) byte {
	if n < len(state.src) {
		return state.src[n]
	}
	return iEOS
}

// peek returns the next byte without advancing the position
func (state *state) peek() byte {
	if state.pos+1 < len(state.src) {
		return state.src[state.pos+1]
	}
	return iEOS
}

// swallow returns the first non-space byte and advances the position
func (state *state) swallow() byte {
	for state.pos+1 < len(state.src) {
		state.pos++
		if !isSpace(state.src[state.pos]) {
			return state.src[state.pos]
		}
	}
	return iEOS
}

// isDigit checks if the byte is a digit from 0 - 9
func isDigit(r byte) bool {
	// return '0' <= r && r <= '9'
	return r >= '0' && r <= '9'
}

// isSpace checks if the byte is empty space
func isSpace(r byte) bool {
	switch r {
	case '\t', '\n', '\v', '\f', '\r', ' ':
		return true
	}
	return false
}
