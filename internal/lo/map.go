package lo

func MapToSlice[K comparable, V any, R any](in map[K]V, transform func(key K, value V) R) []R {
	out := make([]R, 0, len(in))

	for k := range in {
		out = append(out, transform(k, in[k]))
	}

	return out
}
