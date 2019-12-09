package liquidation

import (
	"github.com/mxmCherry/movavg"
	"gonum.org/v1/gonum/stat"
)

// Accumulations is accumulated volume
type Accumulations struct {
	MA *movavg.SMA
}

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
func (p *Accumulations) LiquidationPrice(leverage float64) (buy, sell float64) {
	lastCentralPrice := p.MA.Avg()

	lev2 := 2 * leverage
	return lastCentralPrice * (lev2 - 1) / lev2, lastCentralPrice * (lev2 + 1) / lev2
}
