package service

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"job-scraper.go/internal/core/application"
	"job-scraper.go/internal/models"
	"job-scraper.go/internal/types"
	"time"
)

type ReportService interface {
	GenerateReport([]models.Job, string) (*excelize.File, error)
}

type reportService struct {
	application.Application
}

func newReportService(application application.Application) ReportService {
	return &reportService{
		application,
	}
}

func (r reportService) GenerateReport(jobs []models.Job, userName string) (*excelize.File, error) {
	f := excelize.NewFile()
	linkedInSheetName := "LinkedIn Jobs"
	if err := f.SetSheetName("Sheet1", linkedInSheetName); err != nil {
		return nil, err
	}

	headers := []string{"Title", "Company", "Job Location", "Company Link", "Apply Link", "Staff Count", "Applicants Count", "Expiry Date", "Work Type"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		if err := f.SetCellValue(linkedInSheetName, cell, header); err != nil {
			return nil, err
		}
	}

	for rowIdx, job := range jobs {
		value := []interface{}{job.Title, job.Company, job.JobLocation, job.Link, job.ApplyLink, job.StaffCount, job.ApplicantsCount, time.Unix(job.ExpiryDate, 0).Format("2006-01-02"), job.WorkType}

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

	if err := f.SaveAs(fmt.Sprintf("%s_LinkedIn_Jobs_Report.xlsx", userName)); err != nil {
		return nil, err
	}

	return f, nil
}
