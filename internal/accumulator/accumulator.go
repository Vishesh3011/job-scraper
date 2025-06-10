package accumulator

import (
	"encoding/json"
	"fmt"
	"io"
	"job-scraper.go/internal/core/application"
	"job-scraper.go/internal/models"
	"net/http"
)

type Accumulator struct {
	JobResponse string
}

func NewAccumulator(application application.Application, keyword string) (*Accumulator, error) {
	url := fmt.Sprintf("https://www.linkedin.com/voyager/api/voyagerJobsDashJobCards?decorationId=com.linkedin.voyager.dash.deco.jobs.search.JobSearchCardsCollectionLite-87&count=100&q=jobSearch&query=(currentJobId:4210338428,origin:JOBS_HOME_SEARCH_BUTTON,keywords:%s,locationUnion:(geoId:106089960),spellCorrectionEnabled:true)&servedEventEnabled=false&start=0", keyword)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("csrf-token", application.GetConfig().CsrfToken)
	req.Header.Set("cookie", application.GetConfig().Cookie)
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
	fmt.Println(len(jobs.MetaData.JobCardPrefetchedQueries[0].PrefetchJobPostingCardUrns))

	return &Accumulator{JobResponse: string(body)}, nil
}
