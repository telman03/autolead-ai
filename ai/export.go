package ai

import (
	"fmt"
	"os"
	"strings"

	"github.com/jung-kurt/gofpdf"
	"github.com/telman03/autolead-ai/scraper"
)

func sanitizeFileName(name string) string {
	return strings.ReplaceAll(strings.ToLower(strings.TrimSpace(name)), " ", "_")
}

func ExportCoverLetterToText(company, title, letter string) error {
	filename := fmt.Sprintf("output/%s_%s.txt", sanitizeFileName(company), sanitizeFileName(title))
	return os.WriteFile(filename, []byte(letter), 0644)
}

func ExportCoverLetterToPDF(company, title, letter string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 12)

	lines := strings.Split(letter, "\n")
	for _, line := range lines {
		pdf.MultiCell(0, 8, line, "", "L", false)
	}

	filename := fmt.Sprintf("output/%s_%s.pdf", sanitizeFileName(company), sanitizeFileName(title))
	return pdf.OutputFileAndClose(filename)
}

func ExportBothFormats(job scraper.Job, letter string) {
	err1 := ExportCoverLetterToText(job.Company, job.Title, letter)
	err2 := ExportCoverLetterToPDF(job.Company, job.Title, letter)

	if err1 == nil {
		fmt.Println("✅ Saved .txt")
	} else {
		fmt.Println("❌ Failed to save .txt:", err1)
	}
	if err2 == nil {
		fmt.Println("✅ Saved .pdf")
	} else {
		fmt.Println("❌ Failed to save .pdf:", err2)
	}
}
