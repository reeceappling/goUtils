package utils

import (
	"encoding/json"
	"errors"
	"golang.org/x/exp/maps"
)

type Set[T comparable] map[T]struct{}

func (s Set[T]) Add(es ...T) {

	for _, e := range es {
		s[e] = struct{}{}
	}
}

func (s Set[T]) Contains(el T) bool {
	_, ok := s[el]
	return ok
}

func (s Set[T]) Remove(el T) {
	delete(s, el)
}

func (s Set[T]) ToSlice() []T {
	return maps.Keys(s)
}

// SetOf returns a set only containing unique values of items passed in
func SetOf[T comparable](items []T) Set[T] {
	set := Set[T]{}
	set.Add(items...)
	return set
}

func SetFrom[T comparable](ts ...T) Set[T] {
	return SetOf(ts)
}

// TODO: TEST ME
func MapUniqueChildren[T any, U comparable](allItems []T, getChild func(T) U) Set[U] {
	set := make(Set[U], 0)
	for _, item := range allItems {
		set.Add(getChild(item))
	}
	return set
}

func (s *Set[T]) UnmarshalJSON(data []byte) error {
	var res []T

	if err := json.Unmarshal(data, &res); err != nil {

		// TODO: remove this inner Unmarshall when instances of the old Marshall have expired
		var mapResult map[T]any
		if err = json.Unmarshal(data, &mapResult); err != nil {
			return errors.New("could not unmarshall Set as slice of T nor map of T to any")
		} else {
			*s = SetOf(maps.Keys(mapResult))
		}
	} else if res != nil {
		*s = SetOf(res)
	}

	return nil
}

func (s Set[T]) MarshalJSON() ([]byte, error) {
	var res []T
	if s != nil {
		res = s.ToSlice()
	}
	return json.Marshal(&res)
}
