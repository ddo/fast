package main

import (
	"fmt"
	"time"
	"flag"
	"os"

	"github.com/ddo/go-fast"
	"github.com/ddo/go-spin"
)

func main() {
	var kb, mb, gb bool
	flag.BoolVar(&kb,"k", false, "Format output in Kbps")
	flag.BoolVar(&mb,"m", false, "Format output in Mbps")
	flag.BoolVar(&gb,"g", false, "Format output in Gbps")

	flag.Parse()

	if kb && (mb || gb) || (mb && kb) {
		fmt.Println("You may have at most one formating switch. Choose either -k, -m, or -g")
		os.Exit(-1)
	}

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

	go func() {
		for Kbps := range KbpsChan {
			status = format(Kbps,kb,mb,gb)
		}

		fmt.Printf("\r%c[2K -> %s\n", 27, status)
	}()

	err = fastCom.Measure(urls, KbpsChan)
	ticker.Stop()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return
}

func formatGbps(Kbps float64)(string,string,float64) {
	f := "%.2f %s"
	unit := "Gbps"
	value := Kbps/1000000
	return f,unit,value
}
func formatMbps(Kbps float64)(string,string,float64) {
	f := "%.2f %s"
	unit := "Mbps"
	value := Kbps/1000
	return f,unit,value
}
func formatKbps(Kbps float64)(string,string,float64) {
	f := "%.f %s"
	unit := "Kbps"
	value := Kbps
	return f,unit,value
}

func format(Kbps float64,kb bool,mb bool, gb bool) string {
	var value float64
	var unit string
	var f string

	if kb {
		f,unit,value=formatKbps(Kbps)
	}else if mb {
		f,unit,value=formatMbps(Kbps)
	}else if gb {
		f,unit,value=formatGbps(Kbps)
	}else if Kbps > 1000000 { // Gbps
		f,unit,value=formatGbps(Kbps)
	} else if Kbps > 1000 { // Mbps
		f,unit,value=formatMbps(Kbps)
	}else{
		f,unit,value=formatKbps(Kbps)
	}

	return fmt.Sprintf(f, value, unit)
}
