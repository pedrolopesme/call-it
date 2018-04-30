package call

import (
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"fmt"
)

// Print results accord to spec in
// github.com/pedrolopesme/call-it/issues/6
func PrintResults(result Result) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"URL", "STATUS", "TIMES", "ELAPSED", "AVG", "MIN", "MAX", "TOTAL AVG"})
	table.SetAutoFormatHeaders(false)

	totalExecution := formatTime(result.totalExecution)
	avgExecution := formatTime(result.avgExecution)
	minExeuction := formatTime(result.minExecution)
	maxExecution := formatTime(result.maxExecution)

	firstLine := true
	for statusCode, times := range result.status {
		if firstLine {
			table.Append([]string{result.URL.String(), strconv.Itoa(statusCode), strconv.Itoa(times), " ", " ", minExeuction, maxExecution, avgExecution})
			firstLine = false
		} else {
			table.Append([]string{"", strconv.Itoa(statusCode), strconv.Itoa(times), " ", " ", " ", " ", " "})
		}
	}

	table.SetFooter([]string{"ELAPSED " + totalExecution, "", "", "", "", "", "", " "})
	table.Render()
}


func formatTime(time float64) (output string){
	output = fmt.Sprintf("%.2f", time) + "s"
	return
}