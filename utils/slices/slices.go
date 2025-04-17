package slices

import "golang.org/x/exp/constraints"

// Chunk takes a slice of values and a chunkSize, then returns a slice of the values in slices of chunkSize
// ex: Chunk({1,2,3,4,5,6,7,8},3) == {{1,2,3}{4,5,6}{7,8}}
func Chunk[T any](slice []T, chunkSize int) [][]T {
	var chunks [][]T
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize

		// necessary check to avoid slicing beyond
		// slice capacity
		if end > len(slice) {
			end = len(slice)
		}

		chunks = append(chunks, slice[i:end])
	}

	return chunks
}

// Sliding takes an input slice and a groupsize.         ex: sliding({1,2,3,4,5,6,7},3) == {{1,2,3},
// For each value in the slice that has groupsize-1 values after it,                       	{2,3,4},
// add a slice of {data[i], data[i+1],...,data[i+groupsize-1]} to the output               	{3,4,5},
func Sliding[T any](data []T, groupSize int) [][]T { //										{4,5,6},
	out := make([][]T, len(data)+1-groupSize) //											{5,6,7}}
	for i := range out {
		newGroup := make([]T, groupSize)
		for n := range newGroup {
			newGroup[n] = data[i+n]
		}
		out[i] = newGroup
	}
	return out
}

// Map returns a same-sized slice as the input data.
//
// The value of each position in the output slice is the same position in the input slice modified by the mapFunc.
//
// e.x. output[i] == mapFunc(data[i])
func Map[I, O any](data []I, mapFunc func(I) O) []O {
	out := make([]O, len(data))
	for i := range data {
		out[i] = mapFunc(data[i])
	}
	return out
}

// MapToMap // TODO: this
func MapToMap[I any, O comparable, U any](data []I, mapFunc func(I) (O, U)) map[O]U {
	out := map[O]U{}
	for i := range data {
		key, val := mapFunc(data[i])
		out[key] = val
	}
	return out
}

// Filter applies the keep function to each value in the input
// The result is a slice of only the values which returned true
func Filter[T any](in []T, keep func(T) bool) (out []T) {
	for i := range in {
		toCheck := in[i]
		if keep(toCheck) {
			out = append(out, toCheck)
		}
	}
	return
}

// TODO: changeMe
// The result is a slice of only the values which returned true
func FilterInPlace[T any](in []T, keep func(T) bool) (out []T) { // TODO: ensure works
	current := 0
	for i := range in {
		toCheck := in[i]
		if keep(toCheck) {
			in[i] = toCheck
			current++
		}
	}
	return in[:current] // TODO: ensure correct
}

// ScanLeft creates a slice of size len(data)+1, where the first item is init.
// All other values in the output are the result of the provided operation
func ScanLeft[T constraints.Ordered](data []T, init T, op func(T, T) T) []T {
	out := make([]T, len(data)+1)
	out[0] = init
	for i := 0; i < len(data); i++ {
		out[i+1] = op(out[i], data[i])
	}
	return out
}

// Zip takes two slices and "zips" them together into pairs, dropping unpaired values
// for example L={1,2,3} R={4,5,6} becomes {{1,4}{2,5}{3,6}}
// and L={1,2,3} R={4,5,6,7} becomes {{1,4}{2,5}{3,6}}
func Zip[T any](L, R []T) [][2]T {
	size := min(len(L), len(R))
	out := make([][2]T, size)
	for i := range out {
		out[i] = [2]T{L[i], R[i]}
	}
	return out
}

func ReverseOf[T any](s []T) []T {
	out := make([]T, len(s))
	for i, v := range s {
		out[len(s)-1-i] = v
	}
	return out
}

func Unique[T comparable](in []T) []T {
	out := []T{}
	temp := map[T]struct{}{}
	for _, item := range in {
		if _, exists := temp[item]; !exists {
			out = append(out, item)
			temp[item] = struct{}{}
		}
	}
	return out
}

// UniqueInPlace modifies the input! The output may be a different size than the input!
func UniqueInPlace[T comparable](in []T) []T {
	temp, next := map[T]struct{}{}, 0
	for _, item := range in {
		if _, exists := temp[item]; !exists {
			temp[item] = struct{}{}
			in[next] = item
			next++
		} else {
			temp[item] = struct{}{}
		}
	}
	return in[:next] // TODO: ensure ok

}
