package objbuf

/* #include <stdlib.h> */
import "C"
import (
	"unsafe"
)

func MallocNew(size int) NewFunc[unsafe.Pointer] {
	return func() (unsafe.Pointer, error) {
		return C.malloc(C.size_t(size)), nil
	}
}

func NewMallocBuffer(size, capacity, limit int) (buf *Buffer[unsafe.Pointer], free func(), err error) {
	buf = New(MallocNew(size), limit)
	if _, err := buf.Grow(capacity); err != nil {
		return nil, nil, err
	}

	return buf, func() {
		for _, obj := range buf.Objects(true, true) {
			C.free(obj)
		}
	}, nil
}
