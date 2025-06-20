package config

type Config struct {
	*linkedinConfig
	*dbConfig
	*emailConfig
}

func NewConfig() (*Config, error) {
	linkedINConfig, err := newLinkedINConfig()
	if err != nil {
		return nil, err
	}

	dbConfig, err := newDBConfig()
	if err != nil {
		return nil, err
	}

	emailConfig, err := newEmailConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		linkedinConfig: linkedINConfig,
		dbConfig:       dbConfig,
		emailConfig:    emailConfig,
	}, nil
}
