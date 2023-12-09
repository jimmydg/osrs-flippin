package main

import (
	"github.com/kylelemons/godebug/pretty"
	"log"
)

var (
	FlagCash     *int64
	FlagFlipKind *int
)

func main() {
	handleFlags()
	err, slic := analysis()
	if err != nil {
		log.Panic(err.Error())
	}

	for i := 0; i < 5; i++ {
		pretty.Print(slic[i])
	}
}
