package client

type Client struct {
	*GoMailClient
	*LinkedInClient
	*TelegramClient
}
