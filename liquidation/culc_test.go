package liquidation

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"
)

func TestUse(t *testing.T) {
	termLength := 21
	liq := NewAccumulation(termLength)

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

	mean := liq.MA.Avg()

	// 頻出するLiquidation上下価格などから、刈り取られやすいレバレッジを見込む
	// もしくは、取引所単位や商品などで選択する
	// 算出価格は、中央価格帯で蓄積した出来高がleverageで構築され、証拠金50%になったとき精算される価格とした
	// TODO: 証拠金割合も引数で調整可能にする予定
	var leverage float64 = 15
	buy, sell := liq.LiquidationPrice(leverage)
	fmt.Printf("BidPrice: %.2f AskPrice: %.2f\n", buy, sell)
	fmt.Printf("CenterPrice: %.2f, BothSideDiff: %.2f/%.2f\n", mean, sell-mean, mean-buy)

}
