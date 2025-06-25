package types

type TELEGRAM_STATE string

const (
	AWAIT_USER_NAME          TELEGRAM_STATE = "AWAIT_USER_NAME"
	AWAIT_JOB_ROLES                         = "AWAIT_JOB_ROLES"
	AWAIT_GEO_IDS                           = "AWAIT_GEO_IDS"
	AWAIT_EMAIL_NOTIFY                      = "AWAIT_EMAIL_NOTIFY"
	AWAIT_EMAIL                             = "AWAIT_EMAIL"
	AWAIT_UPDATE_PREFERENCES                = "AWAIT_UPDATE_PREFERENCES"
)
