package objbuf

import "math/bits"

// bm represents a bitmap
type bm []uint64

/**
 * newBm creates a new bitmap object with the specified size.
 *
 * Parameters:
 * - size (int): The size of the bm object to create.
 *
 * Returns:
 * - bm: The newly created bm object.
 */
func newBm(size int) bm {
	out := make(bm, index(size-1)+1)

	for i := range out {
		out[i] = ^uint64(0)
	}

	return out
}

/**
 * resize resizes the bm to the specified size if necessary.
 * It checks if the resize is required based on the input size and the current length of bm.
 * If the size is less than 1 or the length of bm is already greater than or equal to the calculated new size,
 * it returns the original bm.
 * Otherwise, it creates a new bm with the specified size, copies the content of bm to the new bm, and returns the new bm.
 *
 * @param size The new size to which the bm should be resized.
 * @return bm The resized bm based on the specified size.
 */
func (b bm) resize(size int) bm {
	// check if resize required
	if size < 1 || len(b) >= index(size-1)+1 {
		return b
	}

	out := newBm(size)
	copy(out, b)
	return out
}

/**
 * set updates the bit mask at the specified index i.
 */
func (b bm) set(i int) { b[index(i)] |= mask(i) }

/**
 * unset unsets a specific bit in the bitset.
 *
 * Parameters:
 *     i (int): the index of the bit to unset
 */
func (b bm) unset(i int) { b[index(i)] &^= mask(i) }

/**
 * get returns a boolean value based on the index provided.
 */
func (b bm) get(i int) bool { return b[index(i)]&mask(i) != 0 }

/**
 * first returns the index of the first non-zero element in the bm object.
 * If no non-zero element is found, it returns -1.
 */
func (b bm) first() int {
	for i, u := range b {
		if u != 0 {
			return i<<6 + bits.TrailingZeros64(u)
		}
	}

	return -1
}

/**
 * index calculates the index by shifting the input integer i to the right by 6 bits.
 */
func index(i int) int { return i >> 6 }

/**
 * mask generates a bitmask for the given integer.
 * It shifts the value 1 to the left by the result of the bitwise AND operation between the input integer and 63.
 *
 * @param i The integer for which the bitmask is generated.
 * @return The bitmask as a uint64 value.
 */
func mask(i int) uint64 { return 1 << uint(i&63) }
