package main

import (
	"fmt"
	"github.com/miczone/interval-tree/pkg/interval"
)

func main() {
	intervals := interval.CreateIntervalsFromCsvFile("./data/dhl.csv")
	intervals = intervals[:6]
	var root *interval.IntervalNode
	root = root.BuildIntervalTree(intervals)
	root.PrintIntervalNode()

	interval_search := interval.Interval{Low:2, High: 2, Data: 0}
	result := root.OverlapSearch(&interval_search)

	if result == nil {
		fmt.Println("\nNo overlapping interval")
	} else {
		fmt.Printf("\nOverlaps with low %v, high %v, data %v", result.Low, result.High, result.Data)
	}
}
