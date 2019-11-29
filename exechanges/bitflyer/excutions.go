package bitflyer

import (
	"errors"
	"strings"

	"gonum.org/v1/gonum/mat"

	"gonum.org/v1/gonum/stat"

	"github.com/go-numb/go-bitflyer/v1/public/executions"
)

type Bitflyer struct {
	// Leverage is max levarage
	Leverage int

	Upper *mat.Dense
	Lower *mat.Dense
}

// NewBitflyer is client struct
// col[price, vol, label]
// row:0 上昇によりRekt発生、発生をrow:0に設定、下値目安をrow:2に設定
// row:1 上値と下値の平均を設定
// row:2 下落によりRekt発生、発生をrow:2に設定、上値目安をrow:0に設定
func NewBitlyer(lev int) *Bitflyer {
	return &Bitflyer{
		Leverage: lev,
		Upper:    mat.NewDense(1, 3, nil),
		Lower:    mat.NewDense(1, 3, nil),
	}
}

// Set guide price by rekt liquidation order
// 不利（強制精算）出来高を積み立てる
func (p *Bitflyer) Set(isUpper bool, x interface{}) (midPrice float64, err error) {
	data, ok := x.(executions.Execution)
	if !ok {
		return midPrice, errors.New("does not use data")
	}

	if isUpper { // 上昇により売り建玉の精算が発生
		lowerPrice := data.Price * float64(p.Leverage-1) / float64(p.Leverage)
		midPrice = stat.Mean([]float64{data.Price, lowerPrice}, nil)

		add := mat.NewDense(1, 3, []float64{data.Price, data.Size, 0})
		r, c := p.Upper.Dims()
		stack := mat.NewDense(r+1, c, nil)
		stack.Stack(p.Upper, add)
		p.Upper = stack
	} else { // 下落により買い建玉の精算が発生
		upperPrice := data.Price * float64(p.Leverage+1) / float64(p.Leverage)
		midPrice = stat.Mean([]float64{data.Price, upperPrice}, nil)

		add := mat.NewDense(1, 3, []float64{data.Price, data.Size, 0})
		r, c := p.Lower.Dims()
		stack := mat.NewDense(r+1, c, nil)
		stack.Stack(p.Lower, add)
		p.Lower = stack
	}

	return midPrice, nil
}

func (p *Bitflyer) LTP(price float64) float64 {
	return price
}

func (p *Bitflyer) Volume(vol float64) float64 {
	return vol
}

func (p *Bitflyer) ProspectBandwidth(hasAvgPrice, size float64) (mid, ranges float64) {
	if p.Upper.At(0, 0) == p.Lower.At(0, 0) || // 必要なデータが無ければreturn
		p.Upper.At(0, 0) < p.Lower.At(0, 0) || // 下値が上値より高い場合はreturn
		0 == p.Lower.At(0, 0) { // 下値が0値の場合はreturn
		return mid, ranges
	}

	// 上下不利約定の出来高加重平均を取る
	prices := append(p.Upper.RawRowView(0), p.Lower.RawRowView(0)...)
	volumes := append(p.Upper.RawRowView(1), p.Lower.RawRowView(1)...)

	// Denseの構造上、lengthが違うことはありえない
	// 便宜上のchecks length
	if len(prices) != len(volumes) {
		return mid, ranges
	}

	mid = stat.Mean(prices, volumes)

	// TODO:
	// 段階的な建玉の処分
	// 建玉から最適段階と配分を求める
	// diffCenterPrice := math.Abs(upper or lower - mid)
	// bins :=

	// 中央価格前後で利確
	// TODO:
	// 予想到達価格帯を返し、その価格範囲で建玉処分を出力
	return mid, 0
}

// IsRekt finds id, doesnot usually id
func (p *Bitflyer) IsRekt(e executions.Execution) (price, volume float64, isRekt bool) {
	if !strings.HasPrefix(e.BuyChildOrderAcceptanceID, "JRF") {
		return e.Price, e.Size, true
	}
	if !strings.HasPrefix(e.SellChildOrderAcceptanceID, "JRF") {
		return e.Price, e.Size, true
	}
	return e.Price, e.Size, false
}
