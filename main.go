package main

import (
	"fmt"
	"os"
	"time"

	vegeta "github.com/tsenart/vegeta/lib"
)

func main() {
	rate := uint64(1000) // per second
	duration := 4 * time.Second
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "GET",
		URL:    "http://13.114.102.27:8080/sql.md",
	})
	attacker := vegeta.NewAttacker()
	var report vegeta.Results
	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration) {
		metrics.Add(res)
		report.Add(res)
	}
	metrics.Close()
	report.Close()
	fmt.Printf("Requests: %d\n", metrics.Requests)
	fmt.Printf("Max: %s\n", metrics.Latencies.Max)
	fmt.Printf("Mean: %s\n", metrics.Latencies.Mean)
	fmt.Printf("95th percentile: %s\n", metrics.Latencies.P95)

	f, err := os.OpenFile("index.html", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = vegeta.NewTextReporter(&metrics).Report(os.Stdout)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = vegeta.NewPlotReporter("test", &report).Report(f)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
