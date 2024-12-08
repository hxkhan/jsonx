package rt

import "unsafe"

type tflag uint8
type nameOff int32 // offset to a name
type typeOff int32 // offset to an *rtype

const (
	//	u := &(*tUncommon)(unsafe.Pointer(t)).u
	tflagUncommon tflag = 1 << 0

	// tflagExtraStar means the name in the str field has an
	// extraneous '*' prefix. This is because for most types T in
	// a program, the type *T also exists and reusing the str data
	// saves binary size.
	tflagExtraStar tflag = 1 << 1

	// tflagNamed means the type has a name.
	tflagNamed tflag = 1 << 2

	// tflagRegularMemory means that equal and hash functions can treat
	// this type as a single region of t.size bytes.
	tflagRegularMemory tflag = 1 << 3
)

// Type is the runtime type of go values, must be kept in sync with ../runtime/type.go > _type.
type Type struct {
	Size       uintptr // number of bytes a value of this type takes
	Ptrdata    uintptr // number of bytes in the type that can contain pointers
	Hash       uint32  // hash of type; avoids computation in hash tables
	Tflag      tflag   // extra type information flags
	Align      uint8   // alignment of variable with this type
	FieldAlign uint8   // alignment of struct field with this type
	Kind       uint8   // enumeration for C
	// function for comparing objects of this type
	// (ptr to object A, ptr to object B) -> ==?
	Equal     func(unsafe.Pointer, unsafe.Pointer) bool
	GcData    *byte   // garbage collection data
	Str       nameOff // string form
	PtrToThis typeOff // type for pointer to this type, may be zero
}

// Any is the internal structure representing a value boxed in an empty interface
type Any struct {
	Type unsafe.Pointer // a pointer to rt.Type
	Data unsafe.Pointer // a pointer to the actual data
}

// ArrayType is in place of Any.Type if the kind is an Array
type ArrayType struct {
	Type
	ElemsType unsafe.Pointer // a pointer to rt.Type
	SliceType unsafe.Pointer // a pointer to rt.Type for a slice of this type, incase you ever needed it
	Len       uintptr
}

// SliceData is in place of Any.Data if the kind is a Slice
type SliceData struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}

// SliceType is in place of Any.Type if the kind is a Slice
type SliceType struct {
	Type
	ElemsType unsafe.Pointer // a pointer to an rt.Type
}

// StringData is in place of Any.Data if the kind is a String
type StringData struct {
	Data unsafe.Pointer
	Len  int
}

// PointerType is in place of Any.Type if the kind is a Pointer
type PointerType struct {
	Type
	PointedToType unsafe.Pointer // a pointer to an rt.Type of the value this pointer points to
}

// MapType is in place of Any.Type if the kind is a Map
type MapType struct {
	Type
	KeysType  unsafe.Pointer // a pointer to an rt.Type
	ElemsType unsafe.Pointer // a pointer to an rt.Type
	Bucket    unsafe.Pointer // internal bucket structure

	// function for hashing keys (ptr to key, seed) -> hash
	Hasher     func(unsafe.Pointer, uintptr) uintptr
	KeySize    uint8  // size of key slot
	ValueSize  uint8  // size of value slot
	BucketSize uint16 // size of bucket
	Flags      uint32
}
