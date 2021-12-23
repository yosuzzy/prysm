package slice

import (
	"strings"

	types "github.com/prysmaticlabs/eth2-types"
)

// Subset returns true if the first array is
// completely contained in the second array with time
// complexity of approximately o(n).
func Subset[T uint64|int64](a, b []T) bool {
	if len(a) > len(b) {
		return false
	}

	set := make(map[uint64]uint64, len(b))
	for _, v := range b {
		set[uint64(v)]++
	}

	for _, v := range a {
		if count, found := set[uint64(v)]; !found {
			return false
		} else if count < 1 {
			return false
		} else {
			set[uint64(v)] = count - 1
		}
	}
	return true
}

// Intersection of any number of uint64 slices with time
// complexity of approximately O(n) leveraging a map to
// check for element existence off by a constant factor
// of underlying map efficiency.
func Intersection[T uint64|int64](s ...[]T) []T {
	if len(s) == 0 {
		return []T{}
	}
	if len(s) == 1 {
		return s[0]
	}
	intersect := make([]T, 0)
	m := make(map[T]int)
	for _, k := range s[0] {
		m[k] = 1
	}
	for i, num := 1, len(s); i < num; i++ {
		for _, k := range s[i] {
			// Increment and check only if item is present in both, and no increment has happened yet.
			if _, found := m[k]; found && i == m[k] {
				m[k]++
				if m[k] == num {
					intersect = append(intersect, k)
				}
			}
		}
	}
	return intersect
}

// Union of any number of uint64 slices with time
// complexity of approximately O(n) leveraging a map to
// check for element existence off by a constant factor
// of underlying map efficiency.
func Union[T uint64|int64](s ...[]T) []T {
	if len(s) == 0 {
		return []T{}
	}
	if len(s) == 1 {
		return s[0]
	}
	set := s[0]
	m := make(map[T]bool)
	for i := 1; i < len(s); i++ {
		a := s[i-1]
		b := s[i]
		for j := 0; j < len(a); j++ {
			m[a[j]] = true
		}
		for j := 0; j < len(b); j++ {
			if _, found := m[b[j]]; !found {
				set = append(set, b[j])
			}
		}
	}
	return set
}

// SetUint64 returns a slice with only unique
// values from the provided list of indices.
func SetUint64(a []uint64) []uint64 {
	// Remove duplicates indices.
	intMap := map[uint64]bool{}
	cleanedIndices := make([]uint64, 0, len(a))
	for _, idx := range a {
		if intMap[idx] {
			continue
		}
		intMap[idx] = true
		cleanedIndices = append(cleanedIndices, idx)
	}
	return cleanedIndices
}

// IsUint64Sorted verifies if a uint64 slice is sorted in ascending order.
func IsUint64Sorted(a []uint64) bool {
	if len(a) == 0 || len(a) == 1 {
		return true
	}
	for i := 1; i < len(a); i++ {
		if a[i-1] > a[i] {
			return false
		}
	}
	return true
}

// Not returns the uint64 in slice b that are
// not in slice a with time complexity of approximately
// O(n) leveraging a map to check for element existence
// off by a constant factor of underlying map efficiency.
func Not[T uint64|int64](a, b []T) []T {
	set := make([]T, 0)
	m := make(map[T]bool)

	for i := 0; i < len(a); i++ {
		m[a[i]] = true
	}
	for i := 0; i < len(b); i++ {
		if _, found := m[b[i]]; !found {
			set = append(set, b[i])
		}
	}
	return set
}

// IsIn returns true if a is in b and False otherwise.
func IsIn[T uint64|int64](a T, b []T) bool {
	for _, v := range b {
		if a == v {
			return true
		}
	}
	return false
}

// UnionByteSlices returns the all elements between sets of byte slices.
func UnionByteSlices(s ...[][]byte) [][]byte {
	if len(s) == 0 {
		return [][]byte{}
	}
	if len(s) == 1 {
		return s[0]
	}
	set := s[0]
	m := make(map[string]bool)
	for i := 1; i < len(s); i++ {
		for j := 0; j < len(s[i-1]); j++ {
			m[string(s[i-1][j])] = true
		}
		for j := 0; j < len(s[i]); j++ {
			if _, found := m[string(s[i][j])]; !found {
				set = append(set, s[i][j])
			}
		}
	}
	return set
}

// IntersectionByteSlices returns the common elements between sets of byte slices.
func IntersectionByteSlices(s ...[][]byte) [][]byte {
	if len(s) == 0 {
		return [][]byte{}
	}
	if len(s) == 1 {
		return s[0]
	}
	inter := make([][]byte, 0)
	m := make(map[string]int)
	for _, k := range s[0] {
		m[string(k)] = 1
	}
	for i, num := 1, len(s); i < num; i++ {
		for _, k := range s[i] {
			if _, found := m[string(k)]; found && i == m[string(k)] {
				m[string(k)]++
				if m[string(k)] == num {
					inter = append(inter, k)
				}
			}
		}
	}
	return inter
}

// SplitCommaSeparated values from the list. Example: []string{"a,b", "c,d"} becomes []string{"a", "b", "c", "d"}.
func SplitCommaSeparated(arr []string) []string {
	var result []string
	for _, val := range arr {
		result = append(result, strings.Split(val, ",")...)
	}
	return result
}

// SplitOffset returns the start index of a given list splits into chunks,
// it computes (listsize * index) / chunks.
//
// Spec pseudocode definition:
// def get_split_offset(list_size: int, chunks: int, index: int) -> int:
//     """
//     Returns a value such that for a list L, chunk count k and index i,
//     split(L, k)[i] == L[get_split_offset(len(L), k, i): get_split_offset(len(L), k, i+1)]
//     """
//     return (list_size * index) // chunks
func SplitOffset(listSize, chunks, index uint64) uint64 {
	return (listSize * index) / chunks
}

// IntersectionSlot of any number of types.Slot slices with time
// complexity of approximately O(n) leveraging a map to
// check for element existence off by a constant factor
// of underlying map efficiency.
func IntersectionSlot(s ...[]types.Slot) []types.Slot {
	if len(s) == 0 {
		return []types.Slot{}
	}
	if len(s) == 1 {
		return s[0]
	}
	intersect := make([]types.Slot, 0)
	m := make(map[types.Slot]int)
	for _, k := range s[0] {
		m[k] = 1
	}
	for i, num := 1, len(s); i < num; i++ {
		for _, k := range s[i] {
			// Increment and check only if item is present in both, and no increment has happened yet.
			if _, found := m[k]; found && i == m[k] {
				m[k]++
				if m[k] == num {
					intersect = append(intersect, k)
				}
			}
		}
	}
	return intersect
}

// NotSlot returns the types.Slot in slice b that are
// not in slice a with time complexity of approximately
// O(n) leveraging a map to check for element existence
// off by a constant factor of underlying map efficiency.
func NotSlot(a, b []types.Slot) []types.Slot {
	set := make([]types.Slot, 0)
	m := make(map[types.Slot]bool)

	for i := 0; i < len(a); i++ {
		m[a[i]] = true
	}
	for i := 0; i < len(b); i++ {
		if _, found := m[b[i]]; !found {
			set = append(set, b[i])
		}
	}
	return set
}

// IsInSlots returns true if a is in b and False otherwise.
func IsInSlots(a types.Slot, b []types.Slot) bool {
	for _, v := range b {
		if a == v {
			return true
		}
	}
	return false
}
