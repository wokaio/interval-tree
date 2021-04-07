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
	low float64
	high float64
	data float64
}

type IntervalNode struct {
	interval *Interval
	max float64
	left *IntervalNode
	right *IntervalNode
}

func NewIntervalNode(interval *Interval) *IntervalNode {
	return &IntervalNode{
		interval: interval,
		max: interval.high,
		left: nil,
		right: nil,
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

func main() {
	file, err := os.Open("./dhl.csv")
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
			low: low,
			high: high,
			data: data,
		}

		intervals = append(intervals, interval)
	}

	intervals = intervals[:6]

	lenOfIntervals := len(intervals)
	midIndex := int(lenOfIntervals/2)
	var root *IntervalNode
	tmp := ""
	for i := 0; i < lenOfIntervals; i++ {
		if i%2 == 0 {
			midIndex -= i
			root = root.Insert(intervals[midIndex])
		} else {
			midIndex +=i
			root = root.Insert(intervals[midIndex])
		}

		tmp += strconv.Itoa(midIndex) + ","
	}

	println("tmp:",tmp)


	// root.PrintIntervalNode()

	interval_search := Interval{2.5, 2.5, 0}

	result := root.overlapSearch(&interval_search);

	if result == nil {
		fmt.Println("\nNo overlapping interval")
	} else {
		fmt.Printf("\nOverlaps with low %v, high %v, data %v", result.low, result.high, result.data)
	}
}
