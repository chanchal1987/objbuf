package objbuf

import (
	"strconv"
	"testing"
)

func printBm(t *testing.T, bm bm) {
	t.Helper()

	for i, u := range bm {
		t.Logf("%d: %b", i, u)
	}
}

func shouldPanic(t *testing.T, f func()) {
	t.Helper()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic")
		}
	}()

	f()
}

func TestBm(t *testing.T) {
	t.Run("panic", func(t *testing.T) {
		t.Parallel()

		b := newBm(0)

		shouldPanic(t, func() { b.unset(64) })
		shouldPanic(t, func() { b.set(64) })
		shouldPanic(t, func() { b.get(64) })

		if f := b.first(); f != -1 {
			t.Errorf("expected -1 but got %d", f)
		}
	})

	t.Run("panic 2", func(t *testing.T) {
		t.Parallel()

		b := newBm(64)

		shouldPanic(t, func() { b.unset(64) })
		shouldPanic(t, func() { b.set(64) })
		shouldPanic(t, func() { b.get(64) })

		if f := b.first(); f != 0 {
			t.Errorf("expected 0 but got %d", f)
		}
	})

	tests := []struct{ size, l int }{
		{0, 0},
		{1, 1},
		{64, 1},
		{65, 2},
		{128, 2},
		{129, 3},
		{1024, 16},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(strconv.Itoa(tt.size), func(t *testing.T) {
			t.Parallel()

			b := newBm(tt.size)
			if len(b) != tt.l {
				t.Errorf("expected length %d, got %d", tt.l, len(b))
			}

			for _, u := range b {
				if u != ^uint64(0) {
					t.Errorf("expected all bits to be set")
					printBm(t, b)
				}
			}

			if tt.l != 0 && b.first() != 0 {
				t.Errorf("expected 0, got %d", b.first())
			}

			t.Run("unset", func(t *testing.T) {
				for i := 0; i < tt.l; i++ {
					b.unset(i)

					if b.get(i) {
						t.Errorf("expected bit %d to be unset", i)
						printBm(t, b)
						t.FailNow()
					}

					t.Run("check set bits", func(t *testing.T) {
						for j := i + 1; j < tt.l; j++ {
							if !b.get(j) {
								t.Errorf("expected bit %d to be set", j)
								printBm(t, b)
								t.FailNow()
							}
						}
					})

					t.Run("first", func(t *testing.T) {
						if got := b.first(); got != i+1 {
							t.Errorf("expected %d, got %d", i+1, got)
						}
					})
				}
			})

			t.Run("set", func(t *testing.T) {
				for i := range b {
					b[i] = 0
				}

				for i := 0; i < tt.l; i++ {
					b.set(i)
					if !b.get(i) {
						t.Errorf("expected bit %d to be set", i)
						printBm(t, b)
						t.FailNow()
					}

					t.Run("check unset bits", func(t *testing.T) {
						for j := i + 1; j < tt.l; j++ {
							if b.get(j) {
								t.Errorf("expected bit %d to be unset", j)
								printBm(t, b)
								t.FailNow()
							}
						}
					})
				}
			})
		})
	}
}
