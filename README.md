# go-rage-undertake
Im undertaking executions for you.

Use with my package for diff exchanges price
[go-diff-exchanges](https://github.com/go-numb/go-diff-exchanges)


# Usage 
``` golang
package main

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/go-numb/go-rage-undertake/liquidation"
)

func TestUse(t *testing.T) {
	termLength := 21
	liq := liquidation.NewAccumulation(termLength)

	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	for i := 0; i < termLength; i++ {
		var prices, volumes []float64
		for j := 0; j < 10; j++ {
			var price, volume float64
			if j%2 == 0 { // 10000ドルを中心に100ドル上下する価格
				price = 10000.0 - float64(r.Intn(100))
			} else {
				price = 10000.0 + float64(r.Intn(100))
			}
			volume = math.Exp(rand.NormFloat64())

			// 出現価格と紐づく出来高を取得
			// 出来高は捨ててもOK
			prices = append(prices, price)
			volumes = append(volumes, volume)
		}

		// Set 中央価格算出用に出現価格群と出来高群を保持する
		// 中央価格帯で滞留する出来高を凸として、レバレッジごとのリスク範囲を計算するため
		liq.Set(prices, volumes)
	}

	// Setした価格と出来高から算出した中央価格
	mean := liq.MA.Avg()

	// 頻出するLiquidation上下価格などから、刈り取られやすいレバレッジを見込む
	// もしくは、取引所単位や商品などで選択する
	var leverage float64 = 15

	// 算出価格は、中央価格帯で蓄積した出来高がleverageで構築され、証拠金50%になったとき精算される価格とした
	// TODO: 証拠金割合も引数で調整可能にする予定

	/*
		# cut ratio (ex.
			- cutRatio := 0.0 // (Margin or Collateral) Zero cut Exchange
			- cutRatio := 0.5 // (Margin or Collateral) 50% cut Exchange
			- cutRatio := 0.75 // (Margin or Collateral) 75% cut Exchange
	*/
	cutRatio := 0.0 // (Margin or Collateral) Zero cut Exchange
	buy, sell := liq.LiquidationPrice(leverage, cutRatio)
	fmt.Printf("BidPrice: %.2f AskPrice: %.2f\n", buy, sell)
	fmt.Printf("CenterPrice: %.2f, BothSideDiff: %.2f/%.2f\n", mean, sell-mean, mean-buy)

}


```




## Author
[@numbP](https://twitter.com/_numbp)