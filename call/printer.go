package call

import (
	"fmt"
	"strconv"
)

// Print results formatted by Status
func PrintResults(results map[int]int) {
	printHeader()
	for statusCode, times := range results {
		printLine(statusCode, times)
	}
	printFooter();
}

func printHeader() {
	fmt.Println("\n\n\n")
	fmt.Println("+---------------------------------------------+")
	fmt.Println("+ Status Code         | Times                 +")
	fmt.Println("+---------------------------------------------+")
}

func printFooter() {
	fmt.Println("+---------------------------------------------+")
}

func printLine(statusCode int, times int) {
	fmt.Println("+ " +
		strconv.Itoa(statusCode) +
		"                 | " +
		strconv.Itoa(times) +
		"                    +")
}
