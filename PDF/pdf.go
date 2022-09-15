package PDF

import (
	"fmt"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
)

func _PDFPage() *gofpdf.Fpdf {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	return pdf
}

func _SetHeader(pdf *gofpdf.Fpdf) {
	pdf.Image("GL.jpg", 155, 5, 50, 20, false, "", 0, "")
	pdf.SetAlpha(1.0, "Normal") // images := func() {
	pdf.SetFont("Times", "B", 14)
	pdf.CellFormat(40, 5, time.Now().Format("Mon Jan 2, 2006"), "0", 0, "R", false, 0, "")
	// pdf.Ln(5)
	pdf.SetFont("Times", "BI", 14)
	pdf.CellFormat(143, 30, "Education", "0", 0, "R", false, 0, "")

	pdf.Ln(20)
	pdf.SetFont("Times", "B", 28)
	pdf.SetTextColor(255, 0, 0)
	pdf.CellFormat(130, 5, "Daily Report\n", "0", 0, "R", false, 0, "")
}

func _SetFooter(pdf *gofpdf.Fpdf) {
	pdf.SetFooterFunc(func() {
		// Position at 1.5 cm from bottom
		pdf.SetY(-15)
		// Arial italic 8
		pdf.SetFont("Arial", "I", 8)
		// Text color in gray
		pdf.SetTextColor(128, 128, 128)
		// Page number
		pdf.CellFormat(0, 10, fmt.Sprintf("Page %d", pdf.PageNo()),
			"", 0, "C", false, 0, "")
	})
}

func _SetData(pdf *gofpdf.Fpdf, tableName []string, crud []string, res []string, allCol []string) {
	var (
		idx = 0
	)
	for i := range tableName {
		pdf.SetFont("Times", "B", 20)
		pdf.SetTextColor(70, 95, 195)
		pdf.CellFormat(130, 5, "Information about the '"+strings.ToUpper(tableName[i])+"' table", "0", 0, "R", false, 0, "")
		pdf.Ln(10)
		pdf.SetFont("Times", "B", 14)
		pdf.SetTextColor(70, 95, 19)
		pdf.CellFormat(40, 5, "ColumnName: '"+strings.ToUpper(allCol[i]), "0", 0, "L", false, 0, "")
		pdf.Ln(10)
		tr := pdf.UnicodeTranslatorFromDescriptor("")
		for idx_crud := range crud {
			if len(res[idx]) == 0 {
				pdf.SetFont("Times", "B", 18)
				pdf.SetTextColor(70, 95, 19)
				pdf.CellFormat(40, 5, crud[idx_crud]+"\n", "0", 0, "L", false, 0, "")
				pdf.SetFont("Times", "B", 14)
				pdf.SetTextColor(0, 0, 0)
				pdf.CellFormat(20, 5, "None of data "+crud[idx_crud], "0", 0, "L", false, 0, "")
				idx_crud++
				idx++
				pdf.Ln(15)
			} else {
				pdf.SetFont("Times", "B", 18)
				pdf.SetTextColor(70, 95, 19)
				pdf.CellFormat(40, 5, crud[idx_crud]+"\n", "0", 0, "L", false, 0, "")
				pdf.Ln(5)
				pdf.SetFont("Times", "B", 14)
				pdf.AddUTF8Font("Times", "", "")
				pdf.SetTextColor(0, 0, 0)
				pdf.Ln(5)
				pdf.MultiCell(250, 5, tr(res[idx]), "0", "0", false)
				pdf.Ln(10)
				idx++
			}
		}
	}
}

func _DrawTable(pdf *gofpdf.Fpdf, value [][]string) {
	var (
		header = []string{"Table Name", "Inserted", "Updated", "Deleted", "Total"}
	)
	// Colors, line width and bold font
	pdf.SetFillColor(128, 212, 255)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetDrawColor(128, 0, 0)
	pdf.SetLineWidth(.3)
	pdf.SetFont("", "B", 0)
	// 	Header
	w := []float64{40, 25, 25, 25, 20}
	wSum := 0.0
	for _, v := range w {
		wSum += v
	}
	left := (210 - wSum) / 2
	pdf.SetX(left)
	for j, str := range header {
		pdf.CellFormat(w[j], 7, str, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)
	// Color and font restoration
	pdf.SetFillColor(224, 235, 255)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("", "", 0)
	// 	Data
	fill := false
	for _, c := range value {
		pdf.SetX(left)
		pdf.CellFormat(w[0], 10, c[0], "LR", 0, "L", fill, 0, "")
		pdf.CellFormat(w[1], 10, c[1], "LR", 0, "C", fill, 0, "")
		pdf.CellFormat(w[2], 10, c[2], "LR", 0, "C", fill, 0, "")
		pdf.CellFormat(w[3], 10, c[3], "LR", 0, "C", fill, 0, "")
		pdf.CellFormat(w[4], 10, c[4], "LR", 0, "C", fill, 0, "")
		pdf.Ln(-1)
		fill = !fill
	}
	pdf.SetX(left)
	pdf.CellFormat(wSum, 0, "", "T", 0, "", false, 0, "")
}

func CreatePDF(total [][]string, tableName []string, crud []string, res []string, image []string, allCol []string) (err error) {
	pdf := _PDFPage()
	_SetHeader(pdf)
	_SetFooter(pdf)

	pdf.Ln(25)
	pdf.SetFont("Times", "B", 14)
	_DrawTable(pdf, total)
	pdf.Ln(65)
	for i, img := range image {
		if img != "" {
			if i == 0 {
				pdf.Image(img, 20, 120, 155, 140, false, "", 0, "")
				pdf.SetAlpha(1.0, "Normal")
				pdf.AddPage()
				} else {
				pdf.Image(img, 20, 60, 155, 140, false, "", 0, "")
				pdf.SetAlpha(1.0, "Normal")
				pdf.AddPage()
			}
		}
	}
	_SetData(pdf, tableName, crud, res, allCol)
	err = pdf.OutputFileAndClose("report.pdf")
	return err
}
