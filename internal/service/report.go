package service

import (
	"github.com/xuri/excelize/v2"
	"job-scraper.go/internal/models"
	"job-scraper.go/internal/types"
	"time"
)

type ReportService interface {
	GenerateReport([]models.Job, string) (*excelize.File, error)
}

type reportService struct{}

func newReportService() ReportService {
	return &reportService{}
}

func (r reportService) GenerateReport(jobs []models.Job, userName string) (*excelize.File, error) {
	f := excelize.NewFile()
	linkedInSheetName := "LinkedIn Jobs"
	if err := f.SetSheetName("Sheet1", linkedInSheetName); err != nil {
		return nil, err
	}

	headers := []string{"Title", "Company", "Job GeoIds", "Company Link", "Apply Link", "Staff Count", "Applicants Count", "Expiry Date", "Work Type"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		if err := f.SetCellValue(linkedInSheetName, cell, header); err != nil {
			return nil, err
		}
	}

	for rowIdx, job := range jobs {
		var expiry string
		if job.ExpiryDate != nil {
			expiry = time.Unix(*job.ExpiryDate, 0).Format("2006-01-02")
		}

		value := []interface{}{job.Title, job.Company, job.JobLocation, job.Link, job.ApplyLink, job.StaffCount, job.ApplicantsCount, expiry, job.WorkType}

		for colIdx, val := range value {
			if job.Platform == types.LINKEDIN {
				cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+2)
				f.SetCellValue(linkedInSheetName, cell, val)
			}
		}
	}

	if err := f.AutoFilter(linkedInSheetName, "A1:I1", []excelize.AutoFilterOptions{}); err != nil {
		return nil, err
	}

	for i := 0; i < len(headers); i++ {
		colLetter, _ := excelize.ColumnNumberToName(i + 1)
		f.SetColWidth(linkedInSheetName, colLetter, colLetter, 20)
	}

	return f, nil
}
