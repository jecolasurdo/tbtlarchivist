package utils

// PanicIfNil will panic if any of the supplied items is nil.
func PanicIfNil(items ...interface{}) {
	for _, i := range items {
		if i == nil {
			panic("a nil argument was supplied")
		}
	}
}
