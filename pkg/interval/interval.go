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
	DeliveryData Delivery
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

	if root.left != nil && interval.Low <= root.left.max{
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
	fmt.Printf("\n{Low: %v, High: %v, DeliveryData: %v}, max: %v", root.interval.Low, root.interval.High, root.interval.DeliveryData, root.max)
	root.right.PrintIntervalNode()
}

func BuildIntervalTree(intervals []Interval) (root *IntervalNode) {
	intervals_len := len(intervals)
	// balance_index := int(intervals_len / 2)
	// tmp := ""

	for i := 0; i < intervals_len; i++ {
		// if i%2 == 0 {
		// 	balance_index -= i
		// } else {
		// 	balance_index += i
		// }

		// if balance_index == intervals_len {
		// 	balance_index = 0
		// }

		// tmp += strconv.Itoa(i) + ","

		root = root.Insert(intervals[i])
	}

	// fmt.Println("tmp %v", tmp)

	return root
}

func (root *IntervalNode) DeliveryCalculator(weight float64, zone string) *Interval {
	interval_search := Interval{Low: weight, High: weight, DeliveryData: Delivery{}}
	var intervals_result, intervals_result_with_zone []Interval

	root.OverlapSearch(&interval_search, &intervals_result)
	
	for _,value := range intervals_result {
		fmt.Printf("\nOverlaps with low %v, high %v, DeliveryData %v", value.Low, value.High, value.DeliveryData)

		if value.DeliveryData.Zone == zone {
			intervals_result_with_zone = append(intervals_result_with_zone, value)
		}
	}

	result_len := len(intervals_result_with_zone)

	if result_len == 0 {
		return nil
	}

	if (result_len == 1) {
		return &intervals_result_with_zone[0]
	}

	min_price := intervals_result_with_zone[0].DeliveryData.Price
	min_price_index := 0

	for key, value := range intervals_result_with_zone {
		if (min_price > value.DeliveryData.Price) {
			min_price = value.DeliveryData.Price
			min_price_index = key
		}
	}

	return &intervals_result_with_zone[min_price_index]
}

func CreateIntervalsFromCsvFile(path string) ([]Interval) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	var escalone = 0.0005
	var distance = 1.0
	var skip_row_title = true
	var column_title = "LB,Zone A,Zone B,Zone C,Zone D,Zone E,Zone F,Zone G,Zone H,Zone I,Zone J,Zone K,Zone L,Zone M,Zone N"
	var map_column_title = strings.Split(column_title, ",")

	read_file := csv.NewReader(file)
	var intervals []Interval
	var low, high, price float64
	var DeliveryData Delivery
	var zone string

	for {
		record, err := read_file.Read()
		if err == io.EOF {
			break
		}

		if skip_row_title == true {
			skip_row_title = false
			continue;
		}

		for key, value := range record {
			if (key == 0) {
				low, err = strconv.ParseFloat(value, 64)
				low = low + escalone
				high = low + distance 
			} else {
				price, err = strconv.ParseFloat(value, 64)
				zone = map_column_title[key]

				DeliveryData = Delivery{
					Zone: zone,
					Price: price,
				}
		
				interval := Interval{
					Low:  low,
					High: high,
					DeliveryData: DeliveryData,
				}
		
				intervals = append(intervals, interval)
			}
		}
	}

	return intervals
}