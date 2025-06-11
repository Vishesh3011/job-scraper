package models

type LinkedInJobCards struct {
	MetaData struct {
		JobCardPrefetchedQueries []struct {
			PrefetchJobPostingCardUrns []string `json:"prefetchJobPostingCardUrns"`
		} `json:"jobCardPrefetchQueries"`
	} `json:"metaData"`
}

type LinkedInJob struct {
	Name           string `json:"name"`
	Title          string `json:"title"`
	ListedAt       int64  `json:"listedAt"`
	JobPostingId   int64  `json:"jobPostingId"`
	CompanyDetails struct {
		WebJobPosting struct {
			CompanyResolutionResult struct {
				Description     string `json:"description"`
				StaffCount      int    `json:"staffCount"`
				StaffCountRange struct {
					Start int `json:"start"`
					End   int `json:"end"`
				} `json:"staffCountRange"`
				Headquarter struct {
					Country        string `json:"country"`
					GeographicArea string `json:"geographicArea"`
					City           string `json:"city"`
					PostalCode     string `json:"postalCode"`
					Line1          string `json:"line1"`
				} `json:"headquarter"`
				Specialities []string `json:"specialities"`
				Url          string   `json:"url"`
				Industries   []string `json:"industries"`
			} `json:"companyResolutionResult"`
		} `json:"com.linkedin.voyager.deco.jobs.web.shared.WebJobPostingCompany"`
	} `json:"companyDetails"`
	JobLocation        string `json:"formattedLocation"`
	CompanyDescription struct {
		Text string `json:"text"`
	} `json:"companyDescription"`
	JobPostingUrl   string `json:"jobPostingUrl"`
	ApplicantsCount int    `json:"applies"`
	ApplyingInfo    struct {
		Closed  bool `json:"closed"`
		Applied bool `json:"applied"`
	} `json:"applyingInfo"`
	ApplyMethod struct {
		ComLinkedinVoyagerJobsOffsiteApply struct {
			ApplyStartersPreferenceVoid bool   `json:"applyStartersPreferenceVoid"`
			CompanyApplyUrl             string `json:"companyApplyUrl"`
			InPageOffsiteApply          bool   `json:"inPageOffsiteApply"`
		} `json:"com.linkedin.voyager.jobs.OffsiteApply"`
	} `json:"applyMethod"`
	ExpireAt                         int64    `json:"expireAt"`
	Country                          string   `json:"country"`
	WorkplaceTypes                   []string `json:"workplaceTypes"`
	FormattedEmploymentStatus        string   `json:"formattedEmploymentStatus"`
	HiringTeamMembersInjectionResult struct {
		HiringTeamMembers []interface{} `json:"hiringTeamMembers"`
	} `json:"allJobHiringTeamMembersInjectionResult"`
}
