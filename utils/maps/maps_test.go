package maps

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMaps(t *testing.T) {
	testMap := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}

	t.Run("MapToSlice", func(t *testing.T) {
		concatPair := func(s string, i int) string { return fmt.Sprintf(`%s%d`, s, i) }
		out := MapToSlice(testMap, concatPair)
		assert.Equal(t, len(testMap), len(out), "size should be consistent")
		for key, val := range testMap {
			assert.Contains(t, out, concatPair(key, val), fmt.Sprintf(`concatenated value for "%s": %d should exist in the output`, key, val))
		}
	})
	t.Run("Remap - MapToMap", func(t *testing.T) {
		// TODO: THIS
	})
}
