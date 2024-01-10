package utils

import (
	"errors"
	"fmt"
	dom "hw-async/domain"
	"os"
	"slices"
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
	if a.Equal(b) {
		return 0
	} else if a.After(b) {
		return 1
	}

	return -1
}

func TimeToString(t time.Time) string {
	h, m, s := t.Clock()

	return fmt.Sprintf("%v:%v:%v", h, m, s)
}

func CountOfRunesInFloat(f float64) int {
	return utf8.RuneCount([]byte(fmt.Sprintf("%.4f", f)))
}

func ClearFile(path string) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return nil // if file does not exist => no need for clear
	}

	return os.Truncate(path, 0)
}

func Min(nums ...float64) float64 {
	return slices.Min(nums)
}

func Max(nums ...float64) float64 {
	return slices.Max(nums)
}
