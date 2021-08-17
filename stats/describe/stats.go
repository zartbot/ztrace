package describe

import (
	"math"
	"sync"
)

var ONLINE_STATS_SYNC_POOL *sync.Pool

//Item is used to store stats value
type Item struct {
	N    float64
	Min  float64
	Max  float64
	Mean float64
	M2   float64
	M3   float64
	M4   float64
}

type Result struct {
	Count float64
	Min   float64
	Max   float64
	Mean  float64
	Std   float64
	Skew  float64
	Kurt  float64
}

func init() {
	ONLINE_STATS_SYNC_POOL = &sync.Pool{New: func() interface{} {
		return &Item{
			N:    0,
			Min:  0,
			Max:  0,
			Mean: 0,
			M2:   0,
			M3:   0,
			M4:   0,
		}
	}}
}

func New() *Item {
	return ONLINE_STATS_SYNC_POOL.Get().(*Item)
}

//Append is used to store new data for stats, highOrder is in range 2~4 for X^2~4 stats
func (i *Item) Append(x float64, highOrder uint8) {
	n1 := i.N
	if n1 == 0 {
		i.Min = x
		i.Max = x
	}
	if x < i.Min {
		i.Min = x
	}
	if x > i.Max {
		i.Max = x
	}
	i.N = i.N + 1
	delta := x - i.Mean
	deltaN := delta / i.N
	deltaN2 := deltaN * deltaN
	term1 := delta * deltaN * n1
	i.Mean = i.Mean + deltaN

	switch highOrder {
	case 4:
		i.M4 = i.M4 + term1*deltaN2*(i.N*i.N-3*i.N+3) + 6*deltaN2*i.M2 - 4*deltaN*i.M3
		i.M3 = i.M3 + term1*deltaN*(i.N-2) - 3*deltaN*i.M2
		i.M2 = i.M2 + term1
	case 3:
		i.M3 = i.M3 + term1*deltaN*(i.N-2) - 3*deltaN*i.M2
		i.M2 = i.M2 + term1
	case 2:
		i.M2 = i.M2 + term1
	default:
	}
}

//Len is return data count
func (i *Item) Len() float64 {
	return i.N
}

//Sum is return data sum
func (i *Item) Sum() float64 {
	return i.Mean * i.N
}

//Variance is return the variance
func (i *Item) Variance() float64 {
	if i.N < 2 {
		return float64(0.0)
	} else {
		return i.M2 / i.N
	}

}

//Std is return the std
func (i *Item) Std() float64 {
	if i.N < 2 {
		return float64(0.0)
	} else {
		return math.Sqrt(i.M2 / i.N)
	}
}

//Skewness is return skewness...
func (i *Item) Skewness() float64 {
	if i.M2 < 1e-14 || i.N <= 3 {
		return float64(0.0)
	} else {
		return math.Sqrt(i.N) * i.M3 / i.M2 / math.Sqrt(i.M2)
	}
}

//Kurtosis is return kurtosis
func (i *Item) Kurtosis() float64 {
	if i.M2 < 1e-14 || i.N <= 4 {
		return float64(0.0)
	} else {
		return (i.N*i.M4)/(i.M2*i.M2) - 3
	}
}
