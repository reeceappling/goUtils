package maps

func MapToSlice[T comparable, U any, V any](in map[T]U, transformer func(T, U) V) (out []V) {
	out = []V{}
	for key, value := range in {
		out = append(out, transformer(key, value))
	}
	return
}

// Remap iterates over all keys and values to produce a new map of keys and values. Keys and values can differ
func Remap[T comparable, U any, V comparable, W any](in map[T]U, transformer func(T, U) (V, W)) (out map[V]W) {
	out = map[V]W{}
	for key, value := range in {
		v, w := transformer(key, value)
		out[v] = w
	}
	return
}
