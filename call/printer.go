package call

import (
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"fmt"
)

// Print results formatted by Status
func PrintResults(result Result) {
	data := [][]string{}
	for statusCode, times := range result.status {
		data = append(data, []string{strconv.Itoa(statusCode), strconv.Itoa(times)})
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Status Code", "Times"})
	table.SetAutoFormatHeaders(false)
	for _, v := range data {
		table.Append(v)
	}
	totalExecution := fmt.Sprintf("%.1f\n", result.totalExecution) + "s"
	table.SetFooter([]string{"Total Execution", totalExecution})

	table.Render()
}
