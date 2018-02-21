package main

import (
	"fmt"
	"time"

	"github.com/ddo/go-fast"
	"github.com/ddo/go-spin"
)

const (
	ERR = "\r internet error. please try again"
)

func main() {
	status := ""
	spinner := spin.New("")

	// output
	ticker := time.NewTicker(100 * time.Millisecond)

	go func() {
		for range ticker.C {
			fmt.Printf("\r %s  %s", spinner.Spin(), status)
		}
	}()
	// output

	fastCom := fast.New()

	// init
	err := fastCom.Init()

	if err != nil {
		fmt.Println(ERR)
		return
	}

	status = "connecting   "

	// get urls
	urls, err := fastCom.GetUrls()

	if err != nil {
		fmt.Println(ERR)
		return
	}

	status = "loading      "

	// measure
	KbpsChan := make(chan float64)

	go func() {
		for Kbps := range KbpsChan {
			status = format(Kbps) + "    "
		}

		fmt.Printf("\r -> %s\n", status)
	}()

	err = fastCom.Measure(urls, KbpsChan)

	if err != nil {
		fmt.Println(ERR)
	}

	ticker.Stop()
	return
}

func format(Kbps float64) string {
	unit := "Kbps"
	f := "%.f %s"

	if Kbps > 1000000 { // Gbps
		f = "%.2f %s"
		unit = "Gbps"
		Kbps /= 1000

	} else if Kbps > 1000 { // Mbps
		f = "%.2f %s"
		unit = "Mbps"
		Kbps /= 1000
	}

	return fmt.Sprintf(f, Kbps, unit)
}
