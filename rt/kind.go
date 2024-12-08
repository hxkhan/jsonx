package rt

import "unsafe"

type Kind uint8

// Keep in sync with reflect
const (
	Nil Kind = iota
	Bool
	Int
	Int8
	Int16
	Int32
	Int64
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	Uintptr
	Float32
	Float64
	Complex64
	Complex128
	Array
	Chan
	Func
	Interface
	Map
	Pointer
	Slice
	String
	Struct
	UnsafePointer

	kindDirectIface = 1 << 5
	kindGCProg      = 1 << 6
	kindMask        = (1 << 5) - 1
)

func (kind Kind) String() string {
	switch kind {
	case Nil:
		return "nil"
	case Bool:
		return "bool"

	case Int:
		return "int"
	case Int8:
		return "int8"
	case Int16:
		return "int16"
	case Int32:
		return "int32"
	case Int64:
		return "int64"

	case Uint:
		return "uint"
	case Uint8:
		return "uint8"
	case Uint16:
		return "uint16"
	case Uint32:
		return "uint32"
	case Uint64:
		return "uint64"
	case Uintptr:
		return "uintptr"

	case Float32:
		return "float32"
	case Float64:
		return "float64"

	case Complex64:
		return "complex64"
	case Complex128:
		return "complex128"

	case Array:
		return "array"
	case Chan:
		return "chan"
	case Func:
		return "func"
	case Interface:
		return "interface"
	case Map:
		return "map"
	case Pointer:
		return "pointer"
	case Slice:
		return "slice"
	case String:
		return "string"
	case Struct:
		return "struct"
	case UnsafePointer:
		return "unsafe.Pointer"
	}
	return "<UNKNOWN>"
}

func ToAny(x *any) *Any {
	return (*Any)(unsafe.Pointer(x))
}

// Of - Returns the kind of a value, can be compared to values like kind.String, kind.Int etc etc..
func KindOf(value any) Kind {
	if value == nil {
		return Nil
	}

	iFace := (*Any)(unsafe.Pointer(&value))
	//iType := (*Type)(iFace.Type)
	//return Kind(iType.Kind & kindMask)
	// This is a lil faster
	return Kind(*(*uint8)(unsafe.Add(iFace.Type, 23)) & kindMask)
}
