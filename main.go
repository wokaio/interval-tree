package main

import (
	"fmt"
	"github.com/miczone/interval-tree/pkg/interval"
)

func main() {
	intervals := interval.CreateIntervalsFromCsvFile("./data/dhl.csv")
	intervals = intervals[:6]
	tree := interval.BuildIntervalTree(intervals)
	tree.PrintIntervalNode()

	interval_search := interval.Interval{Low:2, High: 2, Data: 0}
	var result []interval.Interval
	tree.OverlapSearch(&interval_search, &result)

	for _,value := range result {
		if result == nil {
			fmt.Println("\nNo overlapping interval")
		} else {
			fmt.Printf("\nOverlaps with low %v, high %v, data %v", value.Low, value.High, value.Data)
		}
	}
}
