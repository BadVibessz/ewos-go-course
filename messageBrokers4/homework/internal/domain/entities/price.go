package entities

import (
	"time"
)

type Price struct {
	Ticker string    `json:"ticker"`
	Value  float64   `json:"value"`
	TS     time.Time `json:"ts"`
}
