package slices

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

func TestSlices(t *testing.T) {
	intSlice := []int{1, 2, 3, 4, 5, 6, 7, 8}
	t.Run("Chunk", func(t *testing.T) {
		assert.Equal(t, [][]int{
			{1, 2, 3},
			{4, 5, 6},
			{7, 8}}, Chunk(intSlice, 3))
	})

	t.Run("Sliding", func(t *testing.T) {
		assert.Equal(t, [][]int{
			{1, 2, 3},
			{2, 3, 4},
			{3, 4, 5},
			{4, 5, 6},
			{5, 6, 7},
			{6, 7, 8}}, Sliding(intSlice, 3))

	})

	t.Run("Map", func(t *testing.T) {
		exp := []string{"1", "2", "3", "4", "5", "6", "7", "8"}
		assert.Equal(t, exp, Map(intSlice, strconv.Itoa))
	})

	t.Run("Filter", func(t *testing.T) {
		keepOdd := func(i int) bool { return i%2 == 1 }
		assert.Equal(t, []int{1, 3, 5, 7}, Filter(intSlice, keepOdd))
	})

	t.Run("ScanLeft", func(t *testing.T) {
		sum := func(i, j int) int { return i + j }
		res := ScanLeft(intSlice, -1, sum)
		exp := []int{-1, 0, 2, 5, 9, 14, 20, 27, 35}
		assert.Equal(t, len(intSlice)+1, len(res))
		assert.Equal(t, exp, res)
	})

	t.Run("Zip", func(t *testing.T) {
		a := []int{1, 3, 5, 7, 28}
		b := []int{2, 4, 6, 8}
		exp := [][2]int{{1, 2}, {3, 4}, {5, 6}, {7, 8}}
		assert.Equal(t, exp, Zip(a, b))
	})

	t.Run("ReverseOf", func(t *testing.T) {
		a := []int{1, 3, 5, 7, 28}
		b := ReverseOf(a)
		aRR := ReverseOf(b)
		bRR := ReverseOf(aRR)
		end := len(a) - 1
		assert.Equal(t, len(a), len(b))
		for i := range a {
			require.Equal(t, a[i], b[end-i], fmt.Sprintf(`b should inverse a's %d index`, i))
			require.Equal(t, b[i], aRR[end-i], fmt.Sprintf(`aRR should inverse b's %d index`, i))
		}
		shouldEqual := func(t *testing.T, c, d []int, idx int, cName, dName string) {
			assert.Equal(t, c[idx], d[idx], fmt.Sprintf(`%s[%d] should equal %s[%d]`, cName, idx, dName, idx))
		}
		for i := range a {
			shouldEqual(t, a, aRR, i, "a", "aRR")
			shouldEqual(t, b, bRR, i, "b", "bRR")
		}
	})
	t.Run("Unique", func(t *testing.T) {
		exp := []int{1, 2, 3}
		a := []int{1, 2, 1, 2, 3}
		lenA := len(a)
		act := Unique(a)
		assert.Equal(t, lenA, len(a), "input length should not change")
		assert.Equal(t, len(exp), len(act))
		for i := range exp {
			assert.Equal(t, exp[i], act[i])
		}
	})
	t.Run("UniqueInPlace", func(t *testing.T) {
		exp := []int{1, 2, 3}
		a := []int{1, 2, 1, 2, 3}
		lenA := len(a)
		assert.NotEqual(t, lenA, len(exp), "input length should start different than expected")
		act := UniqueInPlace(a)
		assert.Equal(t, lenA, len(a), "input length should not change")
		assert.NotEqual(t, 1, a[2])
		assert.Equal(t, len(exp), len(act))
		for i, _ := range exp {
			assert.Equal(t, exp[i], act[i])
		}
	})
}
