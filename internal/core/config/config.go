package config

type Config struct {
	*linkedinConfig
	*dbConfig
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

	return &Config{
		linkedinConfig: linkedINConfig,
		dbConfig:       dbConfig,
	}, nil
}
