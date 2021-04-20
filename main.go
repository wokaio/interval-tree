package main

import (
	"fmt"

	fdinterval "github.com/miczone/interval-tree/pkg/interval"
)

func main() {
	intervals, err := fdinterval.CreateIntervalsFromCsvFile("./data/dhl.csv", 0.0005)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	tree := fdinterval.BuildIntervalTree(intervals)
	// tree.PrintIntervalNode()

	var countryIsoCode string = "us"
	intervalPool := fdinterval.NewIntervalPool()
	intervalPool.SetIntervalPtr(countryIsoCode, tree)

	treeItem, err := intervalPool.GetIntervalPtr(countryIsoCode)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	result, err := treeItem.DeliveryCalculatorByZone(2.3, "Zone D")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("\nDeliveryCalculatorByZone")
	fmt.Printf("Low %v, high %v, data %v", result.Low, result.High, result.DeliveryData)

	results, err := treeItem.DeliveryCalculator(2.3)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("\n\nDeliveryCalculator")
	for i := 0; i < len(results); i++ {
		fmt.Printf("Low %v, high %v, data %v\n", results[i].Low, results[i].High, results[i].DeliveryData)
	}
}
