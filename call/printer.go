package call

import (
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
)

// Print results formatted by Status
func PrintResults(results map[int]int) {
	data := [][]string{}
	for statusCode, times := range results {
		data = append(data, []string{strconv.Itoa(statusCode), strconv.Itoa(times)})
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Status Code", "Times"})
	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}
