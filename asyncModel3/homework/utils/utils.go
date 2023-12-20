package utils

import (
	"fmt"
	dom "hw-async/domain"
	"os"
	"strconv"
	"time"
	"unicode/utf8"
)

func CandlePeriodToInt(per dom.CandlePeriod) (int, error) {
	return strconv.Atoi(string(per[:len(per)-1]))
}

func WriteToFile(path string, b []byte, perm os.FileMode) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, perm)
	if err != nil {
		return err
	}

	_, err = f.Write(b)
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	return nil
}

func PriceToCandle(p dom.Price) dom.Candle {
	return dom.Candle{
		Ticker: p.Ticker,
		Open:   p.Value,
		High:   p.Value,
		Low:    p.Value,
		Close:  p.Value,
		TS:     p.TS,
	}
}

func CandleToCsv(c dom.Candle) string {
	return fmt.Sprintf("%s,%s,%f,%f,%f,%f,%s", c.Ticker, c.Period, c.Open, c.High, c.Low, c.Close, c.TS)
}

func TimeComparator(a, b time.Time) int {
	if a == b {
		return 0
	} else if a.After(b) {
		return 1
	} else {
		return -1
	}
}

func TimeToString(t time.Time) string {
	h, m, s := t.Clock()

	return fmt.Sprintf("%v:%v:%v", h, m, s)
}

func CountOfRunesInFloat(f float64) int {
	s := fmt.Sprintf("%.4f", f)

	return utf8.RuneCount([]byte(s))
}
