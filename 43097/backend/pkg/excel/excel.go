package excel

import (
	"hotel-system/internal/dto"
	"strconv"

	"github.com/xuri/excelize/v2"
)

func ExportOccupancyRateSheet(f *excelize.File, data []dto.OccupancyRateResponse) error {
	sheetName := "入住率统计"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}

	headers := []string{"日期", "总房间数", "已入住房间数", "入住率(%)"}
	for i, header := range headers {
		cell := string(rune('A'+i)) + "1"
		f.SetCellValue(sheetName, cell, header)
	}

	style, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#E0E0E0"},
			Pattern: 1,
		},
	})
	if err == nil {
		f.SetCellStyle(sheetName, "A1", "D1", style)
	}

	for i, item := range data {
		row := strconv.Itoa(i + 2)
		f.SetCellValue(sheetName, "A"+row, item.Date)
		f.SetCellValue(sheetName, "B"+row, item.TotalRooms)
		f.SetCellValue(sheetName, "C"+row, item.OccupiedRooms)
		f.SetCellValue(sheetName, "D"+row, item.OccupancyRate)
	}

	f.SetActiveSheet(index)
	return nil
}

func ExportRevenueSheet(f *excelize.File, data []dto.RevenueResponse) error {
	sheetName := "营收统计"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}

	headers := []string{"日期", "房费收入", "其他收入", "总营收", "退房数量"}
	for i, header := range headers {
		cell := string(rune('A'+i)) + "1"
		f.SetCellValue(sheetName, cell, header)
	}

	style, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#E0E0E0"},
			Pattern: 1,
		},
	})
	if err == nil {
		f.SetCellStyle(sheetName, "A1", "E1", style)
	}

	for i, item := range data {
		row := strconv.Itoa(i + 2)
		f.SetCellValue(sheetName, "A"+row, item.Date)
		f.SetCellValue(sheetName, "B"+row, item.RoomRevenue)
		f.SetCellValue(sheetName, "C"+row, item.OtherRevenue)
		f.SetCellValue(sheetName, "D"+row, item.TotalRevenue)
		f.SetCellValue(sheetName, "E"+row, item.CheckOutCount)
	}

	f.SetActiveSheet(index)
	return nil
}

func SaveExcel(f *excelize.File) ([]byte, error) {
	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
