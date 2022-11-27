package cast

import "strconv"

func ToFloat64(val string) float64 {
	float, err := strconv.ParseFloat(val, 64)
	if err != nil {
		panic(err)
	}
	return float
}

func ToUint64(val string) uint64 {
	u, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		panic(err)
	}
	return u
}
