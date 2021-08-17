package quantile

import "errors"

func MapToQuantileStream(m map[string]interface{}, key string) (*Stream, error) {
	if _v, valid := m[key]; valid {
		if v, ok := _v.(*Stream); ok {
			return v, nil
		} else {
			return nil, errors.New("mismatched Type")
		}
	} else {
		return nil, errors.New("Key is not exist")
	}
}

type Result struct {
	Count int
	P50   float64
	P90   float64
	P99   float64
}

func (s *Stream) Result() *Result {
	return &Result{
		Count: s.Count(),
		P50:   s.Query(0.50),
		P90:   s.Query(0.90),
		P99:   s.Query(0.99),
	}
}
