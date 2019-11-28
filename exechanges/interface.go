package exechanges

import "math"

// 各取引所の必要情報を返す
type Exchange interface {
	Set(x interface{}) error

	LTP(price float64) float64
	Volume(vol float64) float64

	// ProspectBandwidth is predict(invers) price and price range band.
	ProspectBandwidth(avgPrice, size float64) (mid, ranges float64)

	IsRekt() (price, volume float64, isRekt bool)
}

func sigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}
