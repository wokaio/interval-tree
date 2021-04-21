package main

import (
	"fmt"

	fdinterval "github.com/miczone/interval-tree/pkg/interval"
)

func main() {
	intervals, min, max, err := fdinterval.CreateIntervalsFromCsvFile("./data/dhl.csv", 0.0005, -1, -1)
	//intervals, min, max, err := fdinterval.CreateIntervalsFromCsvFile("./data/dhl.csv", 0.0005, 4.0, 150.0)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	tree := fdinterval.BuildIntervalTree(intervals, min, max)
	// tree.PrintIntervalNode()

	var countryIsoCode string = "us"
	intervalPool := fdinterval.NewIntervalPool()
	intervalPool.SetIntervalPtr(countryIsoCode, tree)

	treeItem, err := intervalPool.GetIntervalPtr(countryIsoCode)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("\nDeliveryCalculatorByZone")
	result, err := treeItem.DeliveryCalculatorByZone(4.1, "Zone D")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("Low %v, high %v, data %v", result.Low, result.High, result.DeliveryData)

	fmt.Println("\n\nDeliveryCalculator")
	results, err := treeItem.DeliveryCalculator(4.3)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for i := 0; i < len(results); i++ {
		fmt.Printf("Low %v, high %v, data %v\n", results[i].Low, results[i].High, results[i].DeliveryData)
	}
}
