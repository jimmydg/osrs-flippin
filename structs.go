package main

import "time"

type Latest struct {
	Buy      int64 `json:"high"`
	HighTime int64 `json:"highTime"`
	Sell     int64 `json:"low"`
	LowTime  int64 `json:"lowTime"`
}

type Data struct {
	Data map[int64]Latest `json:"data"`
}

type Mapping struct {
	Examine string `json:"examine"`
	Id      int64  `json:"id"`
	Members bool   `json:"members"`
	LowAlch int64  `json:"lowalch"`
	Limit   int64  `json:"limit"`
	Value   int64  `json:"value"`
	HighAlc int64  `json:"highalc"`
	Icon    string `json:"icon"`
	Name    string `json:"name"`
}

type Volume struct {
	Timestamp int64            `json:"timestamp"`
	Data      map[string]int64 `json:"data"`
}

type Analysis struct {
	Name                     string
	Id                       int64
	Flip                     int64
	FlipWithVolume           int64
	FlipPercentageWithVolume float64
	ProfitPotential          float64
	PercentFlip              float64
	High                     int64
	Low                      int64
	TimeHigh                 time.Time
	TimeLow                  time.Time
	TimeSinceLastFlip        time.Time
	TimeSinceLastFlipPretty  string
	Liquidity24hNonce        int64
	Liquidity24hVolume       int64
	Mapping                  *Mapping
}
