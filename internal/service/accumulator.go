package service

import (
	"encoding/json"
	"fmt"
	"io"
	"job-scraper.go/internal/core/application"
	"job-scraper.go/internal/models"
	"job-scraper.go/internal/utils"
	"net/http"
	url3 "net/url"
)

type Accumulator struct {
	JobResponse string
}

func NewAccumulator(application application.Application, user *models.UserInput) (*Accumulator, error) {
	url := fmt.Sprintf("https://www.linkedin.com/voyager/api/voyagerJobsDashJobCards?decorationId=com.linkedin.voyager.dash.deco.jobs.search.JobSearchCardsCollectionLite-87&count=7&q=jobSearch&query=(origin:JOB_SEARCH_PAGE_SEARCH_BUTTON,keywords:%s,locationUnion:(geoId:%s),spellCorrectionEnabled:true)&servedEventEnabled=false&start=0", url3.QueryEscape(user.Keywords[0]), "104769905")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("csrf-token", application.Config().CsrfToken())
	req.Header.Set("cookie", application.Config().Cookie())
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var jobs models.LinkedInJobCards
	if err := json.Unmarshal(body, &jobs); err != nil {
		return nil, err
	}

	var jobIds []string
	for _, str := range jobs.MetaData.JobCardPrefetchedQueries[0].PrefetchJobPostingCardUrns {
		jobId, err := utils.ExtractJobID(str)
		if err != nil {
			return nil, err
		}
		jobIds = append(jobIds, jobId)
	}
	for _, val := range jobIds {
		fmt.Println(val)
	}

	url2 := fmt.Sprintf("https://www.linkedin.com/voyager/api/jobs/jobPostings/%s?decorationId=com.linkedin.voyager.deco.jobs.web.shared.WebFullJobPosting-65&topN=1&topNRequestedFlavors=List(TOP_APPLICANT,IN_NETWORK,COMPANY_RECRUIT,SCHOOL_RECRUIT,HIDDEN_GEM,ACTIVELY_HIRING_COMPANY)", jobIds[0])
	req2, err := http.NewRequest("GET", url2, nil)
	if err != nil {
		return nil, err
	}
	req2.Header.Set("csrf-token", application.Config().CsrfToken())
	req2.Header.Set("cookie", application.Config().Cookie())
	resp2, err := client.Do(req2)

	if err != nil {
		return nil, err
	}
	defer resp2.Body.Close()

	var job models.LinkedInJob
	if err := json.NewDecoder(resp2.Body).Decode(&job); err != nil {
		return nil, err
	}
	fmt.Println(job.JobPostingUrl)
	fmt.Println(job.Title)

	return &Accumulator{JobResponse: string(body)}, nil
}
