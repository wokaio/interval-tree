package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"github.com/miczone/interval-tree/pkg/interval"
)

func main() {
	intervals := interval.CreateIntervalsFromCsvFile("./data/dhl.csv")
	intervals = intervals[:6]
	var root *IntervalNode
	root = root.BuildIntervalTree(intervals)
	root.PrintIntervalNode()

	interval_search := Interval{2, 2, 0}
	result := root.overlapSearch(&interval_search)

	if result == nil {
		fmt.Println("\nNo overlapping interval")
	} else {
		fmt.Printf("\nOverlaps with low %v, high %v, data %v", result.low, result.high, result.data)
	}
}
