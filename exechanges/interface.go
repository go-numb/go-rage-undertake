package exechanges

// 各取引所の必要情報を返す
type Exchange interface {
	Set(x interface{}) error

	LTP(price float64) float64
	Volume(vol float64) float64

	// ProspectBandwidth is predict(invers) price and price range band.
	ProspectBandwidth() (mean, ranges float64)

	BidPrice() float64
	AskPrice() float64

	IsRekt() (price, volume float64, isRekt bool)
}
