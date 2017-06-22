package runtime

// Information from the compiler about the layout of stack frames.
type bitvector struct {
	n        int32 // # of bits
	bytedata *uint8
}
