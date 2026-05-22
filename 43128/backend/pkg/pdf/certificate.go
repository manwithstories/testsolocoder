package pdf

import (
	"fmt"
	"os"
	"time"

	"github.com/signintech/gopdf"
)

func GenerateCertificate(dir, certNo, eventName, itemName, userName string, rank int, score float64) (string, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	fileName := fmt.Sprintf("%s/%s.pdf", dir, certNo)
	p := gopdf.GoPdf{}
	p.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	p.AddPage()
	if err := p.SetFont("Helvetica", "", 24); err != nil {
		return "", err
	}
	p.SetY(100)
	_ = p.Cell(nil, "Certificate of Achievement")
	_ = p.SetFont("Helvetica", "", 14)
	p.SetY(160)
	_ = p.Cell(nil, "Certificate No: "+certNo)
	p.SetY(185)
	_ = p.Cell(nil, "Awarded to: "+userName)
	p.SetY(210)
	_ = p.Cell(nil, "Event: "+eventName)
	p.SetY(235)
	_ = p.Cell(nil, "Item: "+itemName)
	p.SetY(260)
	_ = p.Cell(nil, fmt.Sprintf("Rank: %d", rank))
	p.SetY(285)
	_ = p.Cell(nil, fmt.Sprintf("Score: %.2f", score))
	p.SetY(340)
	_ = p.Cell(nil, "Issued on: "+time.Now().Format("2006-01-02"))
	if err := p.WritePdf(fileName); err != nil {
		return "", err
	}
	return fileName, nil
}
