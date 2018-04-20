package call

import (
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
)

// Print results formatted by Status
func PrintResults(result Result) {
	data := [][]string{}
	for statusCode, times := range result.status {
		data = append(data, []string{strconv.Itoa(statusCode), strconv.Itoa(times)})
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Status Code", "Times"})
	for _, v := range data {
		table.Append(v)
	}
	table.SetFooter([]string{"Total Execution", strconv.FormatFloat(result.totalExecution, 'g', 1, 64) + "s" })
	table.Render()
}
