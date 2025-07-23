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
)

type jobClient struct {
	httpClient *http.Client
}

func newJobClient(hC *http.Client) *jobClient {
	return &jobClient{
		httpClient: hC,
	}
}

func (c *jobClient) GetLinkedInJobIds(geoId, interest, csrfToken, cookie string) ([]string, error) {
	url := fmt.Sprintf("%s?decorationId=com.linkedin.voyager.dash.deco.jobs.search.JobSearchCardsCollectionLite-87&count=%d&q=jobSearch&query=(origin:JOB_SEARCH_PAGE_SEARCH_BUTTON,keywords:%s,locationUnion:(geoId:%s),spellCorrectionEnabled:true)&servedEventEnabled=false&start=0", types.LINKEDIN_JOB_CARDS_FETCH, types.LINKEDIN_JOB_LIMIT, url3.QueryEscape(interest), geoId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("csrf-token", csrfToken)
	req.Header.Set("cookie", cookie)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
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

func (c *jobClient) GetLinkedInJobDetails(jobId, csrfToken, cookie string) (*models.LinkedInJob, error) {
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

//func (c *jobClient) GetGlassdoorJobs(location, interest string) []models.Job {
//	glassdoorUrl := fmt.Sprintf("https://www.glassdoor.com.au/Job/%s-australia-%s-jobs-SRCH_IL.0,16_IC2235932_KO17,23.htm", location, interest)
//	req, err := http.NewRequest("GET", glassdoorUrl, nil)
//	if err != nil {
//		log.Fatalf("error creating request: %v", err)
//	}
//
//	resp, err := c.httpClient.Do(req)
//	if err != nil {
//		log.Fatalf("error executing request: %v", err)
//	}
//	defer resp.Body.Close()
//	body, err := io.ReadAll(resp.Body)
//	if err != nil {
//		log.Fatalf("error reading response body: %v", err)
//	}
//	os.WriteFile("jobs.html", body, 0777)
//
//	doc, err := goquery.NewDocumentFromReader(resp.Body)
//	if err != nil {
//		log.Fatalf("error loading html body: %v", err)
//	}
//
//	var jobs []models.Job
//	doc.Find("li[data-test='jobListing']").Each(func(i int, s *goquery.Selection) {
//		title := s.Find(".JobCard_jobTitle__GLyJ1").Text()
//		link, _ := s.Find(".JobCard_jobTitle__GLyJ1").Attr("href")
//		company := s.Find(".EmployerProfile_compactEmployerName__9MGcV").Text()
//		location := s.Find(".JobCard_location__Ds1fM").Text()
//		description := s.Find(".JobCard_jobDescriptionSnippet__l1tnl > div").First().Text()
//		//skillsRaw := s.Find(".JobCard_jobDescriptionSnippet__l1tnl > div").Eq(1).Text()
//		//listingAge := s.Find(".JobCard_listingAge__jJsuc").Text()
//
//		job := models.ToJobFromGlassdoor(title, company, location, description, link)
//		jobs = append(jobs, *job)
//	})
//
//	return jobs
//}
//
//func (c *jobClient) GetJoraJobs(location, interest string) []models.Job {
//	url := fmt.Sprintf("https://au.jora.com/j?sp=homepage&trigger_source=homepage&q=%s&l=%s", interest, location)
//	req, err := http.NewRequest("GET", url, nil)
//	if err != nil {
//		log.Fatalf("error creating request: %v", err)
//	}
//
//	resp, err := c.httpClient.Do(req)
//	if err != nil {
//		log.Fatalf("error executing request: %v", err)
//	}
//	defer resp.Body.Close()
//
//	body, err := io.ReadAll(resp.Body)
//	if err != nil {
//		log.Fatalf("error reading response body: %v", err)
//	}
//	os.WriteFile("jobs.html", body, 0777)
//
//	doc, err := goquery.NewDocumentFromReader(resp.Body)
//	if err != nil {
//		log.Fatalf("error loading html body: %v", err)
//	}
//
//	var jobs []models.Job
//	doc.Find("div.job-card").Each(func(i int, s *goquery.Selection) {
//		title := s.Find("h2.job-title a").First().Text()
//		company := s.Find("span.job-company").Text()
//		location := s.Find("a.job-location").Text()
//		description := s.Find("div.job-abstract").Text()
//		href, _ := s.Find("h2.job-title a").First().Attr("href")
//		fullURL := "https://au.indeed.com" + href
//
//		job := models.ToJobFromJora(title, company, location, description, fullURL)
//		jobs = append(jobs, *job)
//	})
//
//	return jobs
//}
