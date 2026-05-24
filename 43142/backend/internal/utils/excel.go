package utils

import (
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"
)

type ExcelColumn struct {
	Header string
	Key    string
	Width  float64
}

type ExcelExportConfig struct {
	SheetName string
	Columns   []ExcelColumn
	Data      []map[string]interface{}
	FileName  string
}

func ExportToExcel(config ExcelExportConfig) (string, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheetName := config.SheetName
	if sheetName == "" {
		sheetName = "Sheet1"
	}

	index, err := f.NewSheet(sheetName)
	if err != nil {
		return "", fmt.Errorf("创建工作表失败: %w", err)
	}
	f.SetActiveSheet(index)

	style, err := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 12},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#E0E0E0"}},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "#CCCCCC", Style: 1},
			{Type: "top", Color: "#CCCCCC", Style: 1},
			{Type: "bottom", Color: "#CCCCCC", Style: 1},
			{Type: "right", Color: "#CCCCCC", Style: 1},
		},
	})
	if err != nil {
		return "", fmt.Errorf("创建样式失败: %w", err)
	}

	dataStyle, err := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 11},
		Alignment: &excelize.Alignment{Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "#CCCCCC", Style: 1},
			{Type: "top", Color: "#CCCCCC", Style: 1},
			{Type: "bottom", Color: "#CCCCCC", Style: 1},
			{Type: "right", Color: "#CCCCCC", Style: 1},
		},
	})
	if err != nil {
		return "", fmt.Errorf("创建数据样式失败: %w", err)
	}

	for i, col := range config.Columns {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, col.Header)
		f.SetCellStyle(sheetName, cell, cell, style)
		if col.Width > 0 {
			f.SetColWidth(sheetName, cell, cell, col.Width)
		}
	}

	for rowIdx, rowData := range config.Data {
		for colIdx, col := range config.Columns {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+2)
			if val, ok := rowData[col.Key]; ok {
				f.SetCellValue(sheetName, cell, val)
			}
			f.SetCellStyle(sheetName, cell, cell, dataStyle)
		}
	}

	fileName := config.FileName
	if fileName == "" {
		fileName = fmt.Sprintf("export_%s.xlsx", time.Now().Format("20060102_150405"))
	}

	if err := f.SaveAs(fileName); err != nil {
		return "", fmt.Errorf("保存Excel文件失败: %w", err)
	}

	return fileName, nil
}
