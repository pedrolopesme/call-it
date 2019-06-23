package call

import (
	"fmt"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
)

// PrintResults output results accord to spec in
// github.com/pedrolopesme/call-it/issues/6
func PrintResults(result Result) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"URL", "STATUS", "TIMES", "AVG", "MIN", "MAX", "TOTAL AVG"})
	table.SetAutoFormatHeaders(false)

	totalExecution := formatTime(result.totalExecution)
	avgExecution := formatTime(result.avgExecution)
	minExeuction := formatTime(result.minExecution)
	maxExecution := formatTime(result.maxExecution)

	firstLine := true
	for statusCode, benchmark := range result.status {

		statusAvgExecution := formatTime(benchmark.execution / float64(benchmark.total))

		if firstLine {
			table.Append([]string{
				result.URL.String(),
				strconv.Itoa(statusCode),
				strconv.Itoa(benchmark.total),
				statusAvgExecution,
				minExeuction,
				maxExecution,
				avgExecution})
			firstLine = false
		} else {
			table.Append([]string{
				"",
				strconv.Itoa(statusCode),
				strconv.Itoa(benchmark.total),
				statusAvgExecution,
				" ",
				" ",
				" "})
		}
	}

	table.SetFooter([]string{"ELAPSED " + totalExecution, "", "", "", "", "", " "})
	table.Render()
}

func formatTime(time float64) (output string) {
	output = fmt.Sprintf("%.2f", time) + "s"
	return
}
