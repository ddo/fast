package main

import (
	"fmt"
	"time"

	"github.com/ddo/go-fast"
	"github.com/ddo/go-spin"
)

func main() {
	status := ""
	spinner := spin.New("")

	// output
	ticker := time.NewTicker(100 * time.Millisecond)

	go func() {
		for range ticker.C {
			fmt.Printf("%c[2K %s  %s\r", 27, spinner.Spin(), status)
		}
	}()
	// output

	fastCom := fast.New()

	// init
	err := fastCom.Init()
	if err != nil {
		fmt.Println(err)
		return
	}

	status = "connecting"

	// get urls
	urls, err := fastCom.GetUrls()
	if err != nil {
		fmt.Println(err)
		return
	}

	status = "loading"

	// measure
	KbpsChan := make(chan float64)

	go func() {
		for Kbps := range KbpsChan {
			status = format(Kbps)
		}

		fmt.Printf("\r%c[2K -> %s\n", 27, status)
	}()

	err = fastCom.Measure(urls, KbpsChan)
	if err != nil {
		fmt.Println(err)
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
		Kbps /= 1000000

	} else if Kbps > 1000 { // Mbps
		f = "%.2f %s"
		unit = "Mbps"
		Kbps /= 1000
	}

	return fmt.Sprintf(f, Kbps, unit)
}
