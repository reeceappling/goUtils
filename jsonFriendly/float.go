package jsonFriendly

import (
	"math"
	"strconv"
)

// None of these constants can be "null" because
// a pointer to a Float that is nil and not omitempty will also be unmarshalled to null
const (
	NaNFloatString    = `"NaN"`
	InfPosFloatString = `"+Inf"`
	InfNegFloatString = `"-Inf"`
)

// Float is a float64 that can ALWAYS be marshalled/unmarshalled into JSON.
// It should never error like with NaN or Inf on a normal float64
type Float float64

func (jff Float) MarshalJSON() ([]byte, error) {
	val := float64(jff)
	if math.IsNaN(val) {
		return []byte(NaNFloatString), nil
	}
	if math.IsInf(val, 1) {
		return []byte(InfPosFloatString), nil
	}
	if math.IsInf(val, -1) {
		return []byte(InfNegFloatString), nil
	}
	return []byte(strconv.FormatFloat(val, 'g', -1, 64)), nil
}

func (jff *Float) UnmarshalJSON(input []byte) error {
	str := string(input)
	switch str {
	case NaNFloatString:
		*jff = Float(math.NaN())
	case InfPosFloatString:
		*jff = Float(math.Inf(1))
	case InfNegFloatString:
		*jff = Float(math.Inf(-1))
	default:
		val, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return err
		}
		*jff = Float(val)
	}
	return nil
}
