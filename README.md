# interval-tree
Interval tree

## Installation

```bash
go get -u github.com/miczone/interval-tree
```

## Example (Optional)

```golang
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
```

### Contributing
- Please create an issue in <a href="https://github.com/miczone/interval-tree/issues">issue list</a>.
- Following the golang coding standards. 

### License
The project is under the Apache 2.0 license. See the LICENSE file for details.
