package objbuf_test

import (
	"testing"
	"unsafe"

	"github.com/chanchal1987/objbuf"
)

func BenchmarkMallocBuffer(b *testing.B) {
	const size = 52
	const alpha = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	buf, free, err := objbuf.NewMallocBuffer(size, 10000, 10000)
	if err != nil {
		b.Fatal(err)
	}

	defer free()

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if err := buf.WithObject(func(p unsafe.Pointer) error {
				bytes := (*[size]byte)(p)[:]
				copy(bytes, alpha)
				return nil
			}); err != nil {
				b.Fatal(err)
			}
		}
	})
}
