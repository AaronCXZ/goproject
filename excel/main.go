package main

import (
	"fmt"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func main() {
	f := excelize.NewFile()
	index := f.NewSheet("Sheet2")
	f.SetCellValue("Sheet2", "A2", "DSFGHJK")
	f.SetCellValue("Sheet1", "B2", "ONNJKBNJBJBJHJ")
	f.SetActiveSheet(index)
	err := f.SaveAs("./test.xlsx")
	if err != nil {
		fmt.Println(err)
	}
}
