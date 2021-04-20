package main

import (
	"fmt"
	"github.com/miczone/interval-tree/pkg/interval"
)

func main() {
	intervals := interval.CreateIntervalsFromCsvFile("./data/dhl.csv")
	tree := interval.BuildIntervalTree(intervals)
	// tree.PrintIntervalNode()
	result := tree.DeliveryCalculator(2.12, "Zone D")

	if result == nil {
		fmt.Println("\nNo overlapping interval")
	} else {
		fmt.Printf("\nOverlaps with low %v, high %v, data %v", result.Low, result.High, result.DeliveryData)
	}

	fmt.Println("\nFinished")
}
