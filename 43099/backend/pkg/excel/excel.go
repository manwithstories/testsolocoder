package excel

import (
	"fmt"
	"venue-booking/internal/model"

	"github.com/xuri/excelize/v2"
)

func ExportPayments(payments []model.Payment) (*excelize.File, error) {
	f := excelize.NewFile()
	sheetName := "支付记录"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}

	headers := []string{
		"ID", "订单号", "交易号", "金额", "支付方式", "支付状态",
		"支付时间", "创建时间",
	}

	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	for i, payment := range payments {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), payment.ID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), payment.Order.OrderNo)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), payment.TransactionNo)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), payment.Amount)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), payment.PaymentMethod)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), payment.Status)
		if payment.PaidAt != nil {
			f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), payment.PaidAt.Format("2006-01-02 15:04:05"))
		}
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), payment.CreatedAt.Format("2006-01-02 15:04:05"))
	}

	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	return f, nil
}
