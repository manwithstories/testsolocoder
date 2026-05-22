package excel

import (
	"io"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

type ScoreRow struct {
	UserID   uint
	Score    float64
	TimeUsed string
	Remarks  string
}

func ParseScoreFile(r io.Reader) ([]ScoreRow, error) {
	f, err := excelize.OpenReader(r)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	sheet := f.GetSheetName(0)
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, err
	}
	var result []ScoreRow
	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 2 {
			continue
		}
		uid, _ := strconv.ParseUint(strings.TrimSpace(row[0]), 10, 64)
		sc, _ := strconv.ParseFloat(strings.TrimSpace(row[1]), 64)
		tu := ""
		rm := ""
		if len(row) >= 3 {
			tu = row[2]
		}
		if len(row) >= 4 {
			rm = row[3]
		}
		result = append(result, ScoreRow{UserID: uint(uid), Score: sc, TimeUsed: tu, Remarks: rm})
	}
	return result, nil
}

func GenerateScoreTemplate(w io.Writer) error {
	f := excelize.NewFile()
	sheet := "Sheet1"
	f.SetCellValue(sheet, "A1", "用户ID")
	f.SetCellValue(sheet, "B1", "成绩")
	f.SetCellValue(sheet, "C1", "用时")
	f.SetCellValue(sheet, "D1", "备注")
	return f.Write(w)
}

func ExportStats(w io.Writer, headers []string, rows [][]interface{}) error {
	f := excelize.NewFile()
	sheet := "统计报表"
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}
	for r, row := range rows {
		for c, v := range row {
			cell, _ := excelize.CoordinatesToCellName(c+1, r+2)
			f.SetCellValue(sheet, cell, v)
		}
	}
	return f.Write(w)
}
