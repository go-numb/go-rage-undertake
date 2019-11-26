package bitflyer

import (
	"errors"
	"strings"

	"github.com/gonum/matrix/mat64"

	"github.com/go-numb/go-bitflyer/v1/public/executions"
)

type Bitflyer struct {
	// Leverage is max levarage
	Leverage int

	Prices *mat64.Dense
}

// NewBitflyer is client struct
func NewBitlyer(lev int) *Bitflyer {
	return &Bitflyer{
		Leverage: lev,
		// col[price,vol]
		Prices: mat64.NewDense(0, 2, nil),
	}
}

func (p *Bitflyer) Set(x interface{}) error {
	data, ok := x.(executions.Execution)
	if !ok {
		return errors.New("does not use data")
	}

	row, _ := p.Prices.Dims()
	p.Prices.SetRow(row, []float64{data.Price, data.Size})
	return nil
}

func (p *Bitflyer) LTP(price float64) float64 {
	return price
}

func (p *Bitflyer) Volume(vol float64) float64 {
	return vol
}

// IsRekt finds id, doesnot usually id
func (p *Bitflyer) IsRekt(e executions.Execution) (price, volume float64, isRekt bool) {
	if !strings.HasPrefix(e.BuyChildOrderAcceptanceID, "JRF") {
		return e.Price, e.Size, true
	}
	if !strings.HasPrefix(e.SellChildOrderAcceptanceID, "JRF") {
		return e.Price, e.Size, true
	}
	return p.LTP(e.Price), p.Volume(e.Size), false
}
