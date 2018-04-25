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
	totalExecution := formatTime(result.totalExecution)
	avgExecution := formatTime(result.avgExecution)
	minExeuction := formatTime(result.minExecution)
	maxExecution := formatTime(result.maxExecution)

	table.Append([]string{"-------------------------", "-------------------------"})
	table.Append([]string{"Avg", avgExecution})
	table.Append([]string{"Min", minExeuction})
	table.Append([]string{"Max", maxExecution})
	table.SetFooter([]string{"Total Execution", totalExecution})

	table.Render()
}

func formatTime(time float64) (output string){
	output = fmt.Sprintf("%.2f", time) + "s"
	return
}