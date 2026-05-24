package utils

import (
	"fmt"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type ExcelColumn struct {
	Field string
	Title string
	Width float64
}

func ExportExcel(c *gin.Context, filename string, sheetName string, columns []ExcelColumn, data interface{}) error {
	f := excelize.NewFile()

	if sheetName == "" {
		sheetName = "Sheet1"
	}

	index, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}

	for i, col := range columns {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, col.Title)
		if col.Width > 0 {
			colName, _ := excelize.ColumnNumberToName(i + 1)
			f.SetColWidth(sheetName, colName, colName, col.Width)
		}
	}

	style, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#CCCCCC"}},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	f.SetCellStyle(sheetName, "A1", fmt.Sprintf("%s1", getColumnName(len(columns))), style)

	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Slice {
		return fmt.Errorf("data must be a slice")
	}

	for i := 0; i < v.Len(); i++ {
		item := v.Index(i)
		if item.Kind() == reflect.Ptr {
			item = item.Elem()
		}

		for j, col := range columns {
			cell, _ := excelize.CoordinatesToCellName(j+1, i+2)

			field := item.FieldByName(col.Field)
			if !field.IsValid() {
				continue
			}

			switch field.Kind() {
			case reflect.String:
				f.SetCellValue(sheetName, cell, field.String())
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				f.SetCellValue(sheetName, cell, field.Int())
			case reflect.Float32, reflect.Float64:
				f.SetCellValue(sheetName, cell, field.Float())
			case reflect.Bool:
				f.SetCellValue(sheetName, cell, field.Bool())
			default:
				if field.Type() == reflect.TypeOf(time.Time{}) {
					f.SetCellValue(sheetName, cell, field.Interface().(time.Time).Format("2006-01-02 15:04:05"))
				} else {
					f.SetCellValue(sheetName, cell, fmt.Sprintf("%v", field.Interface()))
				}
			}
		}
	}

	f.SetActiveSheet(index)

	if filename == "" {
		filename = fmt.Sprintf("export_%s.xlsx", time.Now().Format("20060102150405"))
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Transfer-Encoding", "binary")

	return f.Write(c.Writer)
}

func getColumnName(col int) string {
	name, _ := excelize.ColumnNumberToName(col)
	return name
}

func ParseExcel(filePath string, sheetName string) ([]map[string]string, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if sheetName == "" {
		sheetName = f.GetSheetName(0)
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	if len(rows) < 2 {
		return nil, nil
	}

	headers := rows[0]
	var result []map[string]string

	for _, row := range rows[1:] {
		item := make(map[string]string)
		for i, header := range headers {
			if i < len(row) {
				item[header] = row[i]
			}
		}
		result = append(result, item)
	}

	return result, nil
}
