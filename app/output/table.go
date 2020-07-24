package output

import (
	"github.com/olekukonko/tablewriter"
	"os"
)

func Table(tbData []string)  {
	var allData = [][]string{tbData}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Id", "Name", "Pid", "Status", "Date"})
	hColor := tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiCyanColor}
	table.SetHeaderColor(hColor, hColor, hColor, hColor, hColor)

	for _, row := range allData {
		switch row[3] {
		case "running":
			table.Rich(row, []tablewriter.Colors{{tablewriter.FgHiYellowColor}, {tablewriter.FgHiWhiteColor}, {tablewriter.FgHiWhiteColor}, {tablewriter.Bold, tablewriter.FgHiGreenColor}, {tablewriter.FgWhiteColor}})
		case "stopped":
			table.Rich(row, []tablewriter.Colors{{tablewriter.FgHiYellowColor}, {tablewriter.FgHiWhiteColor}, {tablewriter.FgHiWhiteColor}, {tablewriter.Bold, tablewriter.FgYellowColor}, {tablewriter.FgWhiteColor}})
		case "failed":
			table.Rich(row, []tablewriter.Colors{{tablewriter.FgHiYellowColor}, {tablewriter.FgHiWhiteColor}, {tablewriter.FgHiWhiteColor}, {tablewriter.Bold, tablewriter.FgRedColor}, {tablewriter.FgWhiteColor}})
		case "init":
			table.Rich(row, []tablewriter.Colors{{tablewriter.FgHiYellowColor}, {tablewriter.FgHiWhiteColor}, {tablewriter.FgHiWhiteColor}, {tablewriter.Bold}, {tablewriter.FgWhiteColor}})
		}
	}

	table.SetColumnAlignment([]int{tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER})
	table.Render()
}
