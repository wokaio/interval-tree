package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

type Interval struct {
	low  float64
	high float64
	data float64
}

type IntervalNode struct {
	interval *Interval
	max      float64
	left     *IntervalNode
	right    *IntervalNode
}

func NewIntervalNode(interval *Interval) *IntervalNode {
	return &IntervalNode{
		interval: interval,
		max:      interval.high,
		left:     nil,
		right:    nil,
	}
}

func (root *IntervalNode) Insert(interval Interval) *IntervalNode {
	if root == nil {
		return NewIntervalNode(&interval)
	}

	if interval.low < root.interval.low {
		root.left = root.left.Insert(interval)
	} else {
		root.right = root.right.Insert(interval)
	}

	if root.max < interval.high {
		root.max = interval.high
	}

	return root
}

func doOverlap(interval1 *Interval, interval2 *Interval) bool {
	return interval1.low <= interval2.high && interval2.low <= interval1.high
}

func (root *IntervalNode) overlapSearch(interval *Interval) *Interval {
	if root == nil {
		return nil
	}

	if doOverlap(root.interval, interval) {
		return root.interval
	}

	if root.left != nil && root.left.max >= interval.low {
		return root.left.overlapSearch(interval)
	}

	return root.right.overlapSearch(interval)
}

func (root *IntervalNode) PrintIntervalNode() {
	if root == nil {
		return
	}

	root.left.PrintIntervalNode()
	fmt.Printf("\n{low: %v, high: %v}, max: %v", root.interval.low, root.interval.high, root.max)
	root.right.PrintIntervalNode()
}


func (root *IntervalNode) BuildIntervalTree(intervals []Interval) *IntervalNode {
	intervals_len := len(intervals)
	balance_index := int(intervals_len / 2)

	for i := 0; i < intervals_len; i++ {
		if i%2 == 0 {
			balance_index -= i
		} else {
			balance_index += i
		}

		if balance_index == intervals_len {
			balance_index = 0
		}

		root = root.Insert(intervals[balance_index])
	}

	return root
}

func CreateIntervalsFromCsvFile(path string) []Interval {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	read_file := csv.NewReader(file)
	var intervals []Interval
	var low, high, data float64

	for {

		record, err := read_file.Read()
		if err == io.EOF {
			break
		}

		for key, value := range record {
			if key == 0 { 
				low, err = strconv.ParseFloat(value, 64)
				high = low + 1
			}

			if key == 1 {
				data, err = strconv.ParseFloat(value, 64)
			}
		}

		interval := Interval{
			low:  low,
			high: high,
			data: data,
		}

		intervals = append(intervals, interval)
	}

	return intervals
}

func main() {
	intervals := CreateIntervalsFromCsvFile("./data/dhl.csv")
	intervals = intervals[:6]
	var root *IntervalNode
	root = root.BuildIntervalTree(intervals)

	interval_search := Interval{2, 2, 0}
	result := root.overlapSearch(&interval_search)

	if result == nil {
		fmt.Println("\nNo overlapping interval")
	} else {
		fmt.Printf("\nOverlaps with low %v, high %v, data %v", result.low, result.high, result.data)
	}
}
