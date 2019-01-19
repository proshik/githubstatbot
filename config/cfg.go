package config

import "github.com/kelseyhightower/envconfig"

const (
	// SERVICENAME contains a service name prefix which used in ENV variables
	SERVICENAME = "GITHUBSTATBOT"
)

type Config struct {
	Port               string `envconfig:"PORT",default:"8080"`
	StaticFilesDir     string `default:"./static"`
	GitHubClientId     string `required:"true"`
	GitHubClientSecret string `required:"true"`
	TelegramToken      string `required:"true"`
	AuthBasicUsername  string `default:"user"`
	AuthBasicPassword  string `default:"password"`
	DbUrl              string `envconfig:"DATABASE_URL"`
}

func Load() (*Config, error) {

	cfg := new(Config)

	err := envconfig.Process(SERVICENAME, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
