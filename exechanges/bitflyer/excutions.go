package bitflyer

import (
	"errors"
	"math"
	"strings"

	"gonum.org/v1/gonum/stat"

	"github.com/gonum/matrix/mat64"

	"github.com/go-numb/go-bitflyer/v1/public/executions"
)

type Bitflyer struct {
	// Leverage is max levarage
	Leverage int

	Dense *mat64.Dense
}

// NewBitflyer is client struct
// col[price, vol, label]
// row:0 上昇によりRekt発生、発生をrow:0に設定、下値目安をrow:2に設定
// row:1 上値と下値の平均を設定
// row:2 下落によりRekt発生、発生をrow:2に設定、上値目安をrow:0に設定
func NewBitlyer(lev int) *Bitflyer {
	return &Bitflyer{
		Leverage: lev,
		Dense:    mat64.NewDense(3, 3, nil),
	}
}

// Set guide price by rekt liquidation order
func (p *Bitflyer) Set(isRektBuySide bool, x interface{}) error {
	data, ok := x.(executions.Execution)
	if !ok {
		return errors.New("does not use data")
	}

	if !isRektBuySide { // 上昇により売り建玉の精算が発生
		lowerPrice := data.Price * float64(p.Leverage-1) / float64(p.Leverage)
		midPrice := stat.Mean([]float64{data.Price, lowerPrice}, nil)
		p.Dense.SetRow(0, []float64{data.Price, data.Size})
		p.Dense.SetRow(1, []float64{midPrice, 0})
		p.Dense.SetRow(2, []float64{lowerPrice, 0})
	} else { // 下落により買い建玉の精算が発生
		upperPrice := data.Price * float64(p.Leverage+1) / float64(p.Leverage)
		midPrice := stat.Mean([]float64{data.Price, upperPrice}, nil)
		p.Dense.SetRow(0, []float64{upperPrice, 0, 0})
		p.Dense.SetRow(1, []float64{midPrice, 0, 0})
		p.Dense.SetRow(2, []float64{data.Price, data.Size, 0})
	}

	return nil
}

func (p *Bitflyer) LTP(price float64) float64 {
	return price
}

func (p *Bitflyer) Volume(vol float64) float64 {
	return vol
}

func (p *Bitflyer) ProspectBandwidth(avgPrice, size float64) (mid, ranges float64) {
	mid = p.Dense.At(1, 0)
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
