package bitflyer

import (
	"errors"
	"math"
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
func (p *Bitflyer) Set(isRektBuySide bool, x interface{}) (midPrice float64, err error) {
	data, ok := x.(executions.Execution)
	if !ok {
		return midPrice, errors.New("does not use data")
	}

	if !isRektBuySide { // 上昇により売り建玉の精算が発生
		lowerPrice := data.Price * float64(p.Leverage-1) / float64(p.Leverage)
		midPrice = stat.Mean([]float64{data.Price, lowerPrice}, nil)

		r, _ := p.Upper.Dims()
		p.Upper.SetRow(r+1, []float64{data.Price, data.Size})
		r, _ = p.Lower.Dims()
		p.Lower.SetRow(r+1, []float64{lowerPrice, 0})
	} else { // 下落により買い建玉の精算が発生
		upperPrice := data.Price * float64(p.Leverage+1) / float64(p.Leverage)
		midPrice = stat.Mean([]float64{data.Price, upperPrice}, nil)

		r, _ := p.Upper.Dims()
		p.Upper.SetRow(r+1, []float64{upperPrice, 0, 0})
		r, _ = p.Lower.Dims()
		p.Lower.SetRow(r+1, []float64{data.Price, data.Size, 0})
	}

	return midPrice, nil
}

func (p *Bitflyer) LTP(price float64) float64 {
	return price
}

func (p *Bitflyer) Volume(vol float64) float64 {
	return vol
}

func (p *Bitflyer) ProspectBandwidth(avgPrice, size float64) (mid, ranges float64) {
	// p.Upper.
	mid = p.Upper.At(1, 0)
	diff := math.Abs(avgPrice - mid)
	// 段階的な建玉の処分
	// 建玉から最適段階と配分を求める
	// bins := diff

	_ = diff

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
