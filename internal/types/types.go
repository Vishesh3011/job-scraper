package types

const (
	LINKEDIN_JOB_FETCH       string = "https://www.linkedin.com/voyager/api/jobs/jobPostings/"
	LINKEDIN_JOB_CARDS_FETCH        = "https://www.linkedin.com/voyager/api/voyagerJobsDashJobCards"
)

const (
	LINKEDIN_JOB_LIMIT int = 10
)

type Platform string

const (
	LINKEDIN  Platform = "LinkedIn"
	GLASSDOOR Platform = "Glassdoor"
	INDEED    Platform = "Indeed"
	REMOTEOK  Platform = "RemoteOK"
)
