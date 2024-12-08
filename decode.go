package jsonx

import (
	"strconv"
	"unsafe"
)

// Flag - Options for the Decoder
type Flag int8

// Flags
const (
	// FlagUnprocessedNumbers will make sure to leave the numbers as jsonx.Number instead of float64
	FlagUnprocessedNumbers Flag = iota
)

// Decode - Parses and decodes the input. Flags are for specific cases!
func Decode(input []byte, flags ...Flag) (any, error) {
	// Create a default decoder
	dec := state{src: input, pos: -1}

	// Parse flags
	for _, v := range flags {
		if v == FlagUnprocessedNumbers {
			dec.returnFloats = false
		}
	}

	// Compose the structure
	out, err := dec.compose()
	// Just to make sure that out is nothing except for nil on failure
	if err != nil {
		return nil, err
	}
	return out, nil
}

// compose - Returns Objects, Arrays, Strings and Numbers with err nil, rest as actual char but err
func (state *state) compose() (any, error) {
	char := state.swallow()

	switch char {
	case iQuotation:
		// Get string length
		start := state.pos + 1

	COUNT_STR_LEN:
		char = state.peek()
		if char == iEOS {
			// Unfinished strings are illegal, so error out
			return iEOS, ErrDefault
		} else if char != iQuotation {
			state.pos++
			goto COUNT_STR_LEN
		}
		state.pos++ // trailing quotation
		return string(state.src[start:state.pos]), nil

	case iLeftBracket:
		var array Array

	PARSE_ITEM:
		element, err := state.compose()
		if err == nil {
			array = append(array, element)

			// Next should be an iComma or an iRightBracket
			switch state.swallow() {
			case iComma: // There is more to come...
				goto PARSE_ITEM
			case iRightBracket: // The end has been reached.
				return array, nil
			default: // Error; Unexpected Token: Expected a Comma or a RightBracket
				return nil, ErrDefault
			}
		} else if element == iRightBracket && len(array) == 0 {
			// Array had nothing...
			return array, nil
		}
		// Error ... element can be iEOS or something that does not make sense here
		return nil, err

	case iLeftBrace:
		var obj = make(Object)

	PARSE_PAIR:
		// Parse key
		char = state.swallow()
		if char == iQuotation {
			// Get string length
			start := state.pos + 1

		COUNT_KEY_LEN:
			char = state.peek()
			if char == iEOS {
				// Unfinished strings are illegal, so error out
				return iEOS, ErrDefault
			} else if char != iQuotation {
				state.pos++
				goto COUNT_KEY_LEN
			}
			state.pos++ // trailing quotation
			key := unsafe.String(&state.src[start], state.pos-start)

			/*
				The map does not keep a reference to the original string, but rather stores its own copy of the string's data.
				This means the string is safe to use as a map key even if the original string was derived from an unsafe.String conversion,
				provided the memory backing the string remains valid until after the string has been copied into the map.
			*/

			// Swallow a colon
			if state.swallow() == iColon {
				// Compose value
				value, err := state.compose()
				if err == nil {
					obj[key] = value

					// Swallow next... should be an iComma or an iRightBrace
					char = state.swallow()
					if char == iComma {
						// There is more to come...
						goto PARSE_PAIR
					} else if char == iRightBrace {
						// The end has been reached.
						return obj, nil
					}
					// Unexpected Token: Expected a ',' or a '}'
					return nil, ErrDefault
				}
				return value, err
			}
			// Error Unexpected Token: Expected a Colon
			return nil, ErrDefault
		} else if char == iRightBrace && len(obj) == 0 {
			// Object had nothing...
			return obj, nil
		}
		// Error ... next can be iEOS or something that does not make sense here
		return nil, ErrDefault

	case 't':
		if state.pos+3 < len(state.src) {
			if state.src[state.pos+1] == 'r' &&
				state.src[state.pos+2] == 'u' &&
				state.src[state.pos+3] == 'e' {
				state.pos += 3
				return true, nil
			}
		}
	case 'f':
		if state.pos+4 < len(state.src) {
			if state.src[state.pos+1] == 'a' &&
				state.src[state.pos+2] == 'l' &&
				state.src[state.pos+3] == 's' &&
				state.src[state.pos+4] == 'e' {
				state.pos += 4
				return false, nil
			}
		}
	case 'n':
		if state.pos+3 < len(state.src) {
			if state.src[state.pos+1] == 'u' &&
				state.src[state.pos+2] == 'l' &&
				state.src[state.pos+3] == 'l' {
				state.pos += 3
				return nil, nil
			}
		}
	}

	if isDigit(char) || char == iHyphen {
		var numLen int = 1
		var hasDecimals bool = false

		for char = state.get(state.pos + 1); char != iEOS; char = state.get(state.pos + numLen) {
			if isDigit(char) {
				numLen++
				continue
			} else if char == iDot {
				if !hasDecimals {
					numLen++
					hasDecimals = true
					continue
				}
				// Apparently more than one dot were found
				return nil, ErrDefault
			}

			break
		}

		// Done with parsing number
		str := unsafe.String(&state.src[state.pos], numLen)
		// Now update state's position
		// -1 because state.pos starts with the position of the first digit so account for that
		state.pos += numLen - 1

		if state.returnFloats {
			return stof(str)
		}
		// make string unique
		return Number{hasDecimals, string(([]byte)(str))}, nil
	}
	return char, ErrDefault
}

func stof(str string) (any, error) {
	return strconv.ParseFloat(str, 64)
}
