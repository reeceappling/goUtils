package jsonFriendly

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestJsonFriendlyFloats(t *testing.T) {
	// Test setup
	type jffTemp struct {
		Num Float `json:"num"`
	}
	testJff := func(t *testing.T, flt float64) {
		var result jffTemp
		jff := Float(flt)
		obj := jffTemp{Num: jff}
		bs, err := json.Marshal(obj)
		assert.NoError(t, err)
		assert.NoError(t, json.Unmarshal(bs, &result))
		if math.IsNaN(flt) { // NaN comparison is an edge case
			assert.True(t, math.IsNaN(float64(result.Num)))
		} else {
			assert.Equal(t, flt, float64(result.Num))
		}
	}

	// Run Tests
	t.Run("un/marshalling within larger structs", func(t *testing.T) {
		testJff(t, 3.728)
		testJff(t, math.NaN())
		testJff(t, math.Inf(1))
		testJff(t, math.Inf(2))
	})
}
