package utils

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

type ExcelColumn struct {
	Header string
	Key    string
	Width  float64
}

func ExportExcel(filename string, columns []ExcelColumn, data []map[string]interface{}) ([]byte, error) {
	f := excelize.NewFile()
	sheet := "Sheet1"

	for i, col := range columns {
		cell := fmt.Sprintf("%s1", string(rune(65+i)))
		f.SetCellValue(sheet, cell, col.Header)
		if col.Width > 0 {
			colName, _ := excelize.ColumnNumberToName(i + 1)
			f.SetColWidth(sheet, colName, colName, col.Width)
		}
	}

	styleID, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Color: "#FFFFFF",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#4472C4"},
		},
	})
	f.SetCellStyle(sheet, "A1", fmt.Sprintf("%s1", string(rune(65+len(columns)-1))), styleID)

	for row, item := range data {
		for i, col := range columns {
			cell := fmt.Sprintf("%s%d", string(rune(65+i)), row+2)
			if val, ok := item[col.Key]; ok {
				f.SetCellValue(sheet, cell, val)
			}
		}
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
