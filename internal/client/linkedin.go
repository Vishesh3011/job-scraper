package client

import (
	"encoding/json"
	"fmt"
	"io"
	"job-scraper.go/internal/models"
	"job-scraper.go/internal/types"
	"job-scraper.go/internal/utils"
	"net/http"
	url3 "net/url"
	"os"
)

type linkedInClient struct {
}

func newLinkedInClient() *linkedInClient {
	return &linkedInClient{}
}

func (c *linkedInClient) GetLinkedInJobIds(geoId, interest, csrfToken, cookie string) ([]string, error) {
	url := fmt.Sprintf("%s?decorationId=com.linkedin.voyager.dash.deco.jobs.search.JobSearchCardsCollectionLite-87&count=%d&q=jobSearch&query=(origin:JOB_SEARCH_PAGE_SEARCH_BUTTON,keywords:%s,locationUnion:(geoId:%s),spellCorrectionEnabled:true)&servedEventEnabled=false&start=0", types.LINKEDIN_JOB_CARDS_FETCH, types.LINKEDIN_JOB_LIMIT, url3.QueryEscape(interest), geoId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("csrf-token", csrfToken)
	req.Header.Set("cookie", cookie)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}
	if err := utils.CheckIfResponseIsJSON(body); err != nil {
		return nil, err
	}

	var jobCards models.LinkedInJobCards
	os.WriteFile("jobCards.json", body, os.ModePerm)
	if err := json.Unmarshal(body, &jobCards); err != nil {
		return nil, fmt.Errorf("error unmarshalling job cards: %w", err)
	}

	var jobIds []string
	if len(jobCards.MetaData.JobCardPrefetchedQueries) > 0 && len(jobCards.MetaData.JobCardPrefetchedQueries[0].PrefetchJobPostingCardUrns) > 0 {
		for _, str := range jobCards.MetaData.JobCardPrefetchedQueries[0].PrefetchJobPostingCardUrns {
			jobId, err := utils.ExtractJobID(str)
			if err != nil {
				return nil, fmt.Errorf("error extracting job ID: %w", err)
			}
			jobIds = append(jobIds, jobId)
		}
	}
	return jobIds, nil
}

func (c *linkedInClient) GetLinkedInJobDetails(jobId, csrfToken, cookie string) (*models.LinkedInJob, error) {
	url := fmt.Sprintf("%s%s?decorationId=com.linkedin.voyager.deco.jobs.web.shared.WebFullJobPosting-65&topN=1&topNRequestedFlavors=List(TOP_APPLICANT,IN_NETWORK,COMPANY_RECRUIT,SCHOOL_RECRUIT,HIDDEN_GEM,ACTIVELY_HIRING_COMPANY)", types.LINKEDIN_JOB_FETCH, jobId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("csrf-token", csrfToken)
	req.Header.Set("cookie", cookie)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %w", err)
	}

	defer resp.Body.Close()

	var job models.LinkedInJob
	if err := json.NewDecoder(resp.Body).Decode(&job); err != nil {
		return nil, fmt.Errorf("error decoding job details: %w", err)
	}
	return &job, nil
}
