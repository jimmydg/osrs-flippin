package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"time"
)

var (
	LatestUrl  = "https://prices.runescape.wiki/api/v1/osrs/latest"
	MappingUrl = "https://prices.runescape.wiki/api/v1/osrs/mapping"
	VolumeUrl  = "https://prices.runescape.wiki/api/v1/osrs/volumes"
)

func analysis() (error, []*Analysis) {
	var err error

	latestRes, _ := http.Get(LatestUrl)
	mappingRes, _ := http.Get(MappingUrl)
	volumeRes, _ := http.Get(VolumeUrl)

	fmt.Printf("Res %d \n", latestRes.StatusCode)
	fmt.Printf("Res %d \n", mappingRes.StatusCode)
	fmt.Printf("Res %d \n", volumeRes.StatusCode)

	var latest *Data

	err = decodeJson(&latest, latestRes.Body)
	if err != nil {
		return err, nil
	}

	var mappings []*Mapping
	err = decodeJson(&mappings, mappingRes.Body)
	if err != nil {
		return err, nil
	}

	var volume *Volume
	err = decodeJson(&volume, volumeRes.Body)
	if err != nil {
		return err, nil
	}

	blaat := buildAnalysis(mappings, latest, volume)
	slic := filterByUserCash(blaat)

	return nil, handleFlipAndSlice(slic)
}

func decodeJson(T interface{}, body io.ReadCloser) error {
	defer func(body io.ReadCloser) {
		_ = body.Close()
	}(body)
	return json.NewDecoder(body).Decode(&T)
}

func handleFlags() {
	FlagCash = flag.Int64("cash", 1, "Cash in millions")
	FlagFlipKind = flag.Int(
		"type",
		0,
		"Choose:\n1: By potential profit\n2: By potential flip\n3: By percent flip\n4: Volume combined with percent",
	)
	flag.Parse()
}

func handleFlipAndSlice(slic []*Analysis) []*Analysis {
	switch *FlagFlipKind {
	case 1:
		sort.Slice(slic, func(i, j int) bool {
			return slic[i].ProfitPotential < slic[j].ProfitPotential
		})
	case 2:
		sort.Slice(slic, func(i, j int) bool {
			return slic[i].Flip < slic[j].Flip
		})
	case 3:
		sort.Slice(slic, func(i, j int) bool {
			return slic[i].PercentFlip < slic[j].PercentFlip
		})
	case 4:
		sort.Slice(slic, func(i, j int) bool {
			return slic[i].FlipPercentageWithVolume < slic[j].FlipPercentageWithVolume
		})
	default:
		sort.Slice(slic, func(i, j int) bool {
			return slic[i].Name < slic[j].Name
		})
	}

	return slic
}

func filterByUserCash(analysis []*Analysis) []*Analysis {
	cashToMillions := *FlagCash * 1_000_000

	var slic2 []*Analysis
	for _, v := range analysis {
		if v.Low < cashToMillions {
			slic2 = append(slic2, v)
		}
	}

	var slic3 []*Analysis
	for _, v2 := range slic2 {
		if v2.Liquidity24hNonce > 1000 {
			slic3 = append(slic3, v2)
		}
	}

	//var vec4 []*Analysis
	//for _, v3 := range slic3 {
	//	if v3.Flip > 100 {
	//		vec4 = append(vec4, v3)
	//	}
	//}

	return slic3
}

func buildAnalysis(mapping []*Mapping, latest *Data, volume *Volume) []*Analysis {
	var returnAnalysis []*Analysis
	for i, v := range latest.Data {
		if v.Buy == 0 {
			continue
		}
		liq := volume.Data[strconv.FormatInt(i, 10)]
		liqVol := liq * (v.Buy + v.Sell/2)

		diff := v.Sell - v.Buy
		// 1) Potential profit
		profitPotential := float64(diff * liq)

		// 2) Potential flip
		flipCalc := diff / v.Buy

		// 1) Potential profit
		//profitPotential := float64(flipCalc) * 0.0000003 * (float64(liq) * 0.10)
		//profitPotential := float64(flipCalc * liq)

		// 3) Percentage flip
		flippingPercent := float64((diff - v.Buy) * 100)

		flipWithVolume := flipCalc * (int64(float64(liq) * 0.10))
		flipPercentWithVolume := flippingPercent * (float64(liq) * 0.10)

		timeSinceLastFlip := v.HighTime - v.LowTime

		for _, m := range mapping {
			if m.Id == i {
				boi := &Analysis{
					Name:                     m.Name,
					Flip:                     flipCalc,
					FlipWithVolume:           flipWithVolume,
					FlipPercentageWithVolume: flipPercentWithVolume,
					ProfitPotential:          profitPotential,
					PercentFlip:              flippingPercent,
					High:                     v.Buy,
					Low:                      v.Sell,
					TimeHigh:                 getFormattedDate(v.HighTime),
					TimeLow:                  getFormattedDate(v.LowTime),
					TimeSinceLastFlip:        getFormattedDate(timeSinceLastFlip),
					TimeSinceLastFlipPretty:  getFormattedDate(timeSinceLastFlip).Format("15:03:01"),
					Liquidity24hNonce:        liq,
					Liquidity24hVolume:       liqVol,
					Id:                       i,
					Mapping:                  m,
				}
				returnAnalysis = append(returnAnalysis, boi)
			}
		}

	}

	return returnAnalysis
}

func getFormattedDate(timestamp int64) time.Time {
	t := time.Unix(timestamp, 0)
	return t
	//return t.Format(time.RFC822)
}
