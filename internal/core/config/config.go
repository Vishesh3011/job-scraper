package config

type Config struct {
	*LinkedinConfig
}

func NewConfig() (*Config, error) {
	linkedINConfig, err := newLinkedINConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		LinkedinConfig: linkedINConfig,
	}, nil
}
