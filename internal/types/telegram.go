package types

type TELEGRAM_STATE string

const (
	AWAIT_USER_NAME      TELEGRAM_STATE = "AWAIT_USER_NAME"
	AWAIT_JOB_ROLES                     = "AWAIT_JOB_ROLES"
	AWAIT_GEO_IDS                       = "AWAIT_GEO_IDS"
	AWAIT_COOKIE                        = "AWAIT_COOKIE"
	AWAIT_CSRF_TOKEN                    = "AWAIT_CSRF_TOKEN"
	AWAIT_EMAIL_NOTIFY                  = "AWAIT_EMAIL_NOTIFY"
	AWAIT_EMAIL                         = "AWAIT_EMAIL"
	AWAIT_UPDATE_DETAILS                = "AWAIT_UPDATE_DETAILS"

	SEND_REPORT TELEGRAM_STATE = "SEND_REPORT"
	FINISHED    TELEGRAM_STATE = "FINISHED"
)

type BotPrompt string

const (
	PromptWelcome BotPrompt = "👋 Welcome to JobScraper Bot!\n\nThis bot helps you stay updated with the latest job listings from LinkedIn, directly in Telegram and optionally via email. Let's get you started with a quick setup!\n\nTo begin, please enter your name."

	PromptEnterName         BotPrompt = "📝 Please enter your full name:"
	PromptEnterJobRoles     BotPrompt = "💼 What job roles are you interested in?\n(Separate multiple roles with commas)"
	PromptEnterJobLocations BotPrompt = "🌍 Which job location geo-IDs are you targeting?\n(Separate multiple IDs with commas)\n\nE.g., 102713980 for Australia, 102257491 for India"

	PromptEnterLinkedInCookie    BotPrompt = "🔐 Please paste your LinkedIn session cookie:"
	PromptEnterLinkedInCSRFToken BotPrompt = "🛡️ Please enter your LinkedIn CSRF token:"

	PromptAskEmailReportPreference BotPrompt = "📩 Would you like to receive a daily job report via email? (y/n)"
	PromptEnterEmail               BotPrompt = "📧 Great! Please enter your email address:"

	PromptRegistrationSuccess        BotPrompt = "🎉 You're all set! You've been successfully registered.\nWe'll now start sending job reports to your Telegram and email."
	PromptAccountExistsUpdateRequest BotPrompt = "🔁 It looks like you're already registered!\nWould you like to update your details? (y/n)"
	PromptPreferencesUpdated         BotPrompt = "✅ Your preferences have been updated successfully!"
	PromptReportGenerated            BotPrompt = "📬 Your job report has been generated!\nCheck your Telegram and email inbox."

	PromptErrorProcessingRequest BotPrompt = "⚠️ Oops! Something went wrong. Please try again later."
)
