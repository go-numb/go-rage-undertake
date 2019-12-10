package liquidation

import (
	"github.com/mxmCherry/movavg"
	"gonum.org/v1/gonum/stat"
)

// Accumulations is accumulated volume
type Accumulations struct {
	MA *movavg.SMA
}

// NewAccumulation has length data
func NewAccumulation(length int) *Accumulations {
	return &Accumulations{
		MA: movavg.NewSMA(length),
	}
}

// Set take for central price
func (p *Accumulations) Set(prices, volumes []float64) {
	p.MA.Add(stat.Mean(prices, volumes))
}

// LiquidationPrice culc liquidation price for set Leverage
// Cut ratio is exchange rule.
// Ratio 0.5 is do liquidation when (collateral or margin) *0.5
// Ratio 0 is do liquidation when (collateral or margin) *0
func (p *Accumulations) LiquidationPrice(leverage float64, cutRatio float64) (buy, sell float64) {
	lastCentralPrice := p.MA.Avg()

	f := float64(1) - cutRatio
	ratio := float64(1) / f
	lev2 := ratio * leverage
	return lastCentralPrice * (lev2 - 1) / lev2, lastCentralPrice * (lev2 + 1) / lev2
}
