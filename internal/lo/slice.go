package lo

// CountValuesBy counts the number of occurrences of each value in the collection.
// The mapper function is used to extract the value to count.
// The result is a map where the key is the value and the value is the number of occurrences.
func CountValuesBy[T any, U comparable](collection []T, mapper func(item T) U) map[U]int {
	out := make(map[U]int)

	for i := range collection {
		out[mapper(collection[i])]++
	}

	return out
}

// Map manipulates a slice and transforms it to a slice of another type.
func Map[T any, R any](collection []T, convert func(item T, idx int) R) []R {
	out := make([]R, len(collection))

	for i := range collection {
		out[i] = convert(collection[i], i)
	}

	return out
}

// FlatMap manipulates a slice and transforms it to a slice of another type.
func FlatMap[T any, R any](collection []T, convert func(item T, index int) []R) []R {
	out := make([]R, 0, len(collection))

	for i := range collection {
		out = append(out, convert(collection[i], i)...)
	}

	return out
}

// FilterMap returns a slice which obtained after both filtering and mapping using the given callback function.
// The callback function should return two values:
//   - the result of the mapping operation and
//   - whether the result element should be included or not.
func FilterMap[T any, R any](collection []T, callback func(item T, index int) (R, bool)) []R {
	out := make([]R, 0)

	for i := range collection {
		if r, ok := callback(collection[i], i); ok {
			out = append(out, r)
		}
	}

	return out
}

// ContainsBy returns true if predicate function return true.
func ContainsBy[T any](collection []T, predicate func(item T) bool) bool {
	for i := range collection {
		if predicate(collection[i]) {
			return true
		}
	}

	return false
}
