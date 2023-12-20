package graphic

import (
	"fmt"
	"github.com/notEpsilon/go-pair"
	"hw-async/utils"
	"math"
	"slices"
	"sort"
	"strings"
	"time"
)

type TimePrice pair.Pair[time.Time, float64]

type CandleGraphic struct {
	Ticker     string
	TimePrices []TimePrice
}

func NewCandleGraphic(ticker string, timePrices ...TimePrice) *CandleGraphic {
	// sort
	//slices.SortFunc(timePrices, func(a, b TimePrice) int { return utils.TimeComparator(a.First, b.First) })

	return &CandleGraphic{Ticker: ticker, TimePrices: timePrices}
}

func (g *CandleGraphic) GenerateString() string {
	res := fmt.Sprintf("%s\n", g.Ticker)

	prices := make([]float64, 0, len(g.TimePrices))
	timings := make([]time.Time, 0, len(g.TimePrices))

	for _, tp := range g.TimePrices {
		prices = append(prices, tp.Second)
		timings = append(timings, tp.First)
	}

	// sort prices descending
	sort.Sort(sort.Reverse(sort.Float64Slice(prices)))

	// sort timings ascending
	slices.SortFunc(timings, utils.TimeComparator)

	leftOff := 0

	for i := range prices {
		p := prices[i]

		tIdx := slices.IndexFunc(g.TimePrices, func(tp TimePrice) bool { return tp.Second == p }) // todo: utils
		t := g.TimePrices[tIdx].First

		// get corresponding time from timings slice
		tIdx = slices.IndexFunc(timings, func(tim time.Time) bool { return tim == t })

		// todo: calculate gap in spaces and print ***** there

		priceRunesCount := utils.CountOfRunesInFloat(p)
		off := 3

		leftOff = int(math.Max(float64(leftOff), float64(priceRunesCount+off)))

		spaces := strings.Repeat(" ", (tIdx+1)*(priceRunesCount+off))
		res += fmt.Sprintf("%.4f %s %s\n", p, spaces, "*******")
	}

	res += strings.Repeat(" ", leftOff)
	for _, t := range timings {
		res += fmt.Sprintf("%s  \n", utils.TimeToString(t))
	}

	return res
}
