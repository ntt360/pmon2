package output

import (
	"github.com/olekukonko/tablewriter"
	"os"
)

func DescTable(tbData [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"--", "Desc"})
	hColor := tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiCyanColor}
	table.SetHeaderColor(hColor, hColor)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	for index, row := range tbData {
		if index == 0 {
			switch row[1] {
			case "running":
				table.Rich(row, []tablewriter.Colors{{}, {tablewriter.Bold, tablewriter.FgHiGreenColor}})
				break
			case "stopped":
				table.Rich(row, []tablewriter.Colors{{}, {tablewriter.Bold, tablewriter.FgHiYellowColor}})
				break
			case "failed":
				table.Rich(row, []tablewriter.Colors{{}, {tablewriter.Bold, tablewriter.FgRedColor}})
			default:
				table.Append(row)
			}
		} else {
			table.Append(row)
		}
	}

	//table.SetRowLine(true)
	table.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT})
	table.Render()
}

func Table(tbData [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Id", "Name", "Pid", "Status", "Date"})
	hColor := tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiCyanColor}
	table.SetHeaderColor(hColor, hColor, hColor, hColor, hColor)

	for _, row := range tbData {
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
