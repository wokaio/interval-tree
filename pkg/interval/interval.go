package fdinterval

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"reflect"
	"sort"
	"strconv"
)

var (
	floatType  = reflect.TypeOf(float64(0))
	stringType = reflect.TypeOf("")
)

type Interval struct {
	Low          float64
	High         float64
	DeliveryData Delivery
}

type IntervalList []Interval

func (item_list IntervalList) Len() int {
	return len(item_list)
}

func (item_list IntervalList) Less(i, j int) bool {
	return item_list[i].DeliveryData.Price > item_list[j].DeliveryData.Price
}

func (item_list IntervalList) Swap(i, j int) {
	item_list[i], item_list[j] = item_list[j], item_list[i]
}

type Delivery struct {
	Zone  string
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

	if root.left != nil && interval.Low <= root.left.max {
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

	for i := 0; i < intervals_len; i++ {
		root = root.Insert(intervals[i])
	}

	return root
}

func ConvertI2Float(number interface{}) (float64, error) {
	switch i := number.(type) {
	case float64:
		return i, nil
	case float32:
		return float64(i), nil
	case int64:
		return float64(i), nil
	case int32:
		return float64(i), nil
	case int:
		return float64(i), nil
	case uint64:
		return float64(i), nil
	case uint32:
		return float64(i), nil
	case uint:
		return float64(i), nil
	case string:
		return strconv.ParseFloat(i, 64)
	default:
		v := reflect.ValueOf(number)
		v = reflect.Indirect(v)
		if v.Type().ConvertibleTo(floatType) {
			fv := v.Convert(floatType)
			return fv.Float(), nil
		} else if v.Type().ConvertibleTo(stringType) {
			sv := v.Convert(stringType)
			s := sv.String()
			return strconv.ParseFloat(s, 64)
		} else {
			return math.NaN(), fmt.Errorf("Error to convert %v to float64", v.Type())
		}
	}
}

func (root *IntervalNode) DeliveryCalculatorByZone(weight interface{}, zone string) (*Interval, error) {
	weightf, err := ConvertI2Float(weight)
	if err != nil {
		return nil, err
	}
	
	if weightf < 0 {
		return nil, errors.New("Weight must be > 0")
	}

	interval_search := Interval{Low: weightf, High: weightf, DeliveryData: Delivery{}}
	var intervals_result, intervals_result_with_zone []Interval

	root.OverlapSearch(&interval_search, &intervals_result)
	for _, value := range intervals_result {
		if value.DeliveryData.Zone == zone {
			intervals_result_with_zone = append(intervals_result_with_zone, value)
		}
	}

	result_len := len(intervals_result_with_zone)
	if result_len == 0 {
		return nil, errors.New("No overlapping interval")
	}

	if result_len == 1 {
		return &intervals_result_with_zone[0], nil
	}

	sort.Sort(IntervalList(intervals_result_with_zone))
	return &intervals_result_with_zone[0], nil
}

func (root *IntervalNode) DeliveryCalculator(weight interface{}) ([]Interval, error) {
	weightf, err := ConvertI2Float(weight)
	if err != nil {
		return nil, err
	}
	
	if weightf < 0 {
		return nil, errors.New("Weight must be > 0")
	}

	interval_search := Interval{Low: weightf, High: weightf, DeliveryData: Delivery{}}
	var intervals_result, intervals_result_with_zone []Interval

	root.OverlapSearch(&interval_search, &intervals_result)
	for _, value := range intervals_result {
		intervals_result_with_zone = append(intervals_result_with_zone, value)
	}

	result_len := len(intervals_result_with_zone)
	if result_len == 0 {
		return nil, errors.New("No overlapping interval")
	}

	if result_len > 1 {
		sort.Sort(IntervalList(intervals_result_with_zone))
	}

	return intervals_result_with_zone, nil
}

func StringToFloat64(numStr string) (float64, error) {
	return strconv.ParseFloat(numStr, 64)
}

func CreateIntervalsFromCsvFile(path string, step float64, min float64, max float64) ([]Interval, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var distance = 1.0
	var skip_row_title = true
	var map_column_title []string

	read_file := csv.NewReader(file)
	if read_file == nil {
		return nil, errors.New("Can not read CSV file")
	}

	var intervals []Interval
	var low float64 = 0.00
	var high float64 = 0.00
	var price float64 = 0.00
	var DeliveryData Delivery
	var zone string

	for {
		record, err := read_file.Read()
		if err == io.EOF {
			break
		}

		if len(record) == 0 {
			continue
		}

		if skip_row_title == true {
			skip_row_title = false
			map_column_title = record
			continue
		}

		high, err = StringToFloat64(record[0])
		if err != nil {
			continue
		}
		high = high + step
		high = math.Floor(high*10000) / 10000
		low = high - distance
		low = math.Floor(low*10000) / 10000

		if (min > 0) || (max > 0) {
			if (high > max) || (low < min) {
				continue
			}
		}

		for idx := 1; idx < len(record); idx++ {
			price, err = StringToFloat64(record[idx])
			zone = map_column_title[idx]

			DeliveryData = Delivery{
				Zone:  zone,
				Price: price,
			}

			interval := Interval{
				Low:          low,
				High:         high,
				DeliveryData: DeliveryData,
			}

			intervals = append(intervals, interval)
		}
	}

	return intervals, nil
}
