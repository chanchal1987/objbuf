package objbuf

import "sync"

// NewFunc creates a new object of type T
type NewFunc[T any] func() (T, error)

// Buffer is a generic buffer for any type T.
type Buffer[T any] struct {
	cond  sync.Cond  // cond is used to synchronize access to the buffer.
	store []T        // store holds the actual objects in the buffer.
	bm    bm         // bm is a bitmap used to track, which objects are in use.
	lim   int        // lim is the maximum number of objects that can be in use.
	new   NewFunc[T] // new is a function that creates a new object of type T.
}

/**
 * Creates a new Buffer instance with the specified limit and NewFunc function.
 *
 * @param New The function to create a new instance of type T.
 * @param limit The maximum number of elements the buffer can hold.
 *
 * @return A pointer to the newly created Buffer instance.
 */
func New[T any](New NewFunc[T], limit int) *Buffer[T] {
	return &Buffer[T]{
		cond: sync.Cond{L: new(sync.Mutex)},
		lim:  limit,
		new:  New,
	}
}

func (buf *Buffer[T]) growUnsafe(i int) (bool, error) {
	if buf.new == nil {
		return false, nil
	}

	if len(buf.store)+i >= buf.lim {
		return false, nil
	}

	objs := make([]T, i)

	for j := 0; j < i; j++ {
		obj, err := buf.new()
		if err != nil {
			return false, err
		}

		objs[j] = obj
	}

	buf.store = append(buf.store, objs...)
	buf.bm = buf.bm.resize(len(buf.store))
	return true, nil
}

func (buf *Buffer[T]) Grow(i int) (bool, error) {
	buf.cond.L.Lock()
	defer buf.cond.L.Unlock()

	return buf.growUnsafe(i)
}

func (buf *Buffer[T]) Objects(used, unused bool) []T {
	out := make([]T, 0, len(buf.store))

	buf.cond.L.Lock()
	defer buf.cond.L.Unlock()

	for i, c := range buf.store {
		if unused && buf.bm.get(i) {
			out = append(out, c)
		}

		if used && !buf.bm.get(i) {
			out = append(out, c)
		}
	}

	return out[:len(out):len(out)]
}

func (buf *Buffer[T]) WithObject(f func(T) error) error {
	var idx int

	buf.cond.L.Lock()

	for {
		if idx = buf.bm.first(); idx < 0 || idx >= len(buf.store) {
			if ok, err := buf.growUnsafe(1); err != nil {
				buf.cond.L.Unlock()
				return err
			} else if ok {
				idx = len(buf.store) - 1
				break
			}

			buf.cond.Wait()
		} else {
			break
		}
	}

	buf.bm.unset(idx)
	buf.cond.L.Unlock()

	defer func() {
		buf.cond.L.Lock()
		buf.bm.set(idx)
		buf.cond.L.Unlock()
		buf.cond.Signal()
	}()

	return f(buf.store[idx])
}
