package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ddo/go-fast"
	"github.com/ddo/go-spin"
)

func main() {
	var kb, mb, gb, silent bool
	flag.BoolVar(&kb, "k", false, "Format output in Kbps")
	flag.BoolVar(&mb, "m", false, "Format output in Mbps")
	flag.BoolVar(&gb, "g", false, "Format output in Gbps")
	flag.BoolVar(&silent, "silent", false, "Surpress all output except for the final result")

	flag.Parse()

	if kb && (mb || gb) || (mb && kb) {
		fmt.Println("You may have at most one formating switch. Choose either -k, -m, or -g")
		os.Exit(-1)
	}

	status := ""
	spinner := spin.New("")

	// output
	ticker := time.NewTicker(100 * time.Millisecond)
	if !silent {
		go func() {
			for range ticker.C {
				fmt.Printf("%c[2K %s  %s\r", 27, spinner.Spin(), status)
			}
		}()
	}
	// output

	fastCom := fast.New()

	// init
	err := fastCom.Init()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	status = "connecting"

	// get urls
	urls, err := fastCom.GetUrls()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	status = "loading"

	// measure
	KbpsChan := make(chan float64)
	done := make(chan bool)
	go func(done chan bool) {
		var value, units string
		for Kbps := range KbpsChan {
			value, units = format(Kbps, kb, mb, gb)
			// don't print the units of measurement if explicitly asked for
			if kb || mb || gb {
				status = fmt.Sprintf("%s", value)
			} else {
				status = fmt.Sprintf("%s %s", value, units)
			}

		}
		if silent {
			fmt.Printf("%s\n", status)
		} else {
			fmt.Printf("\r%c[2K -> %s\n", 27, status)
		}

		done <- true
	}(done)

	err = fastCom.Measure(urls, KbpsChan)
	ticker.Stop()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// finish reading KbpsChan so that the result always gets printed
	<- done

	return
}

func formatGbps(Kbps float64) (string, string, float64) {
	f := "%.2f"
	unit := "Gbps"
	value := Kbps / 1000000
	return f, unit, value
}
func formatMbps(Kbps float64) (string, string, float64) {
	f := "%.2f"
	unit := "Mbps"
	value := Kbps / 1000
	return f, unit, value
}
func formatKbps(Kbps float64) (string, string, float64) {
	f := "%.f"
	unit := "Kbps"
	value := Kbps
	return f, unit, value
}

func format(Kbps float64, kb bool, mb bool, gb bool) (string, string) {
	var value float64
	var unit string
	var f string

	if kb {
		f, unit, value = formatKbps(Kbps)
	} else if mb {
		f, unit, value = formatMbps(Kbps)
	} else if gb {
		f, unit, value = formatGbps(Kbps)
	} else if Kbps > 1000000 { // Gbps
		f, unit, value = formatGbps(Kbps)
	} else if Kbps > 1000 { // Mbps
		f, unit, value = formatMbps(Kbps)
	} else {
		f, unit, value = formatKbps(Kbps)
	}

	strValue := fmt.Sprintf(f, value)
	return strValue, unit
}
