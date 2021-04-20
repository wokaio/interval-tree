package interval

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type Interval struct {
	Low  float64
	High float64
	Data interface {}
}

type Delivery struct {
	Zone string
	Price float64
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
		max:      interval.High,
		left:     nil,
		right:    nil,
	}
}

func (root *IntervalNode) Insert(interval Interval) *IntervalNode {
	if root == nil {
		return NewIntervalNode(&interval)
	}

	if interval.Low < root.interval.Low {
		root.left = root.left.Insert(interval)
	} else {
		root.right = root.right.Insert(interval)
	}

	if root.max < interval.High {
		root.max = interval.High
	}

	return root
}

func doOverlap(interval1 *Interval, interval2 *Interval) bool {
	return interval1.Low <= interval2.High && interval2.Low <= interval1.High
}

func (root *IntervalNode) OverlapSearch(interval *Interval, intervals *[]Interval) {
	if root == nil {
		return 
	}

	if doOverlap(root.interval, interval) {
		*intervals = append(*intervals, *root.interval)
	}

	if root.left != nil && root.left.max >= interval.Low {
		root.left.OverlapSearch(interval, intervals)
		return 
	}

	root.right.OverlapSearch(interval, intervals)
}

func (root *IntervalNode) PrintIntervalNode() {
	if root == nil {
		return
	}

	root.left.PrintIntervalNode()
	fmt.Printf("\n{Low: %v, High: %v}, max: %v", root.interval.Low, root.interval.High, root.max)
	root.right.PrintIntervalNode()
}

func BuildIntervalTree(intervals []Interval) (root *IntervalNode) {
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

func CreateIntervalsFromCsvFile(path string) ([]Interval) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	var str_zone = "LB,Zone A,Zone B,Zone C,Zone D,Zone E,Zone F,Zone G,Zone H,Zone I,Zone J,Zone K,Zone L,Zone M,Zone N"
	var map_column_zone = strings.Split(str_zone, ",")

	read_file := csv.NewReader(file)
	var intervals []Interval
	var low, high, price float64
	var data interface {}
	var zone string
	var row_number = 1

	for {
		record, err := read_file.Read()
		if err == io.EOF {
			break
		}

		if row_number == 1 {
			row_number ++
			continue;
		}

		for key, value := range record {
			if (key == 0) {
				low, err = strconv.ParseFloat(value, 64)
				high = low + 1
			} else {
				price, err = strconv.ParseFloat(value, 64)
				zone = map_column_zone[key]

				data = Delivery{
					Zone: zone,
					Price: price,
				}
		
				interval := Interval{
					Low:  low,
					High: high,
					Data: data,
				}
		
				intervals = append(intervals, interval)
			}
		}
	}

	return intervals
}