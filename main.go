package main

import (
	"fmt"
	"strings"
)

func printTable(list []int) {
	var header, values strings.Builder

	for idx, value := range list {
		fmt.Fprintf(&header, "%d | ", idx)
		fmt.Fprintf(&values, "%d | ", value)
	}

	splitter := strings.Repeat("-", header.Len()-1)

	fmt.Printf("%s\n%s\n%s\n",header.String(), splitter, values.String())


}

func main() {
	fmt.Println("Hello, World!")

	arr1 := []int{0, 0, 0, 0, 0}
	// arr2 := make([]int, 10)

	printTable(arr1)


}