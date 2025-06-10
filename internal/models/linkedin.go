package models

type LinkedInJobCardPrefetchQuery struct {
	PrefetchJobPostingCardUrns []string `json:"prefetchJobPostingCardUrns"`
}

type LinkedInMetaData struct {
	JobCardPrefetchedQueries []LinkedInJobCardPrefetchQuery `json:"jobCardPrefetchQueries"`
}

type LinkedInJobCards struct {
	MetaData *LinkedInMetaData `json:"metaData"`
}
