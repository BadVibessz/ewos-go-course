package graphic

import (
	"fmt"
	"hw-async/utils"
	"math"
	"slices"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/notEpsilon/go-pair"
)

type TimePrice pair.Pair[time.Time, float64]

type CandleGraphic struct {
	Ticker     string
	TimePrices []TimePrice
}

func New(ticker string, timePrices ...TimePrice) *CandleGraphic {
	return &CandleGraphic{Ticker: ticker, TimePrices: timePrices}
}

func (g *CandleGraphic) GenerateString(piv rune) string {
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
	lastTimeLen := 0

	for i := range prices {
		p := prices[i]

		tIdx := slices.IndexFunc(g.TimePrices, func(tp TimePrice) bool { return tp.Second == p })
		t := g.TimePrices[tIdx].First

		// get corresponding time from timings slice
		tIdx = slices.IndexFunc(timings, func(tim time.Time) bool { return tim == t })

		priceRunesCount := utils.CountOfRunesInFloat(p)
		off := 2

		timeStr := utils.TimeToString(t)
		timeLen := utf8.RuneCountInString(timeStr)

		if lastTimeLen == 0 {
			lastTimeLen = timeLen
		}

		spaces := strings.Repeat(" ", (tIdx)*(off+lastTimeLen))
		pivots := strings.Repeat(string(piv), timeLen)

		res += fmt.Sprintf("%.4f %s %s\n", p, spaces, pivots)

		leftOff = int(math.Max(float64(leftOff), float64(priceRunesCount+off)))
		lastTimeLen = timeLen
	}

	res += strings.Repeat(" ", leftOff)
	for _, t := range timings {
		res += fmt.Sprintf("%s  ", utils.TimeToString(t))
	}

	return res + "\n"
}
