package config

import "github.com/kelseyhightower/envconfig"

const (
	// SERVICENAME contains a service name prefix which used in ENV variables
	SERVICENAME = "GITHUBSTATBOT"
)

type Config struct {
	Port               string `envconfig:"PORT" default:"8080"`
	DbUrl              string `envconfig:"DATABASE_URL" default:"postgres://postgres:password@localhost:5432/githubstatbot?sslmode=disable"`
	GitHubClientId     string `required:"true"`
	GitHubClientSecret string `required:"true"`
	TelegramToken      string `required:"true"`
	StaticFilesDir     string `default:"./static"`
	AuthBasicUsername  string `default:"user"`
	AuthBasicPassword  string `default:"password"`
}

func Load() (*Config, error) {
	cfg := new(Config)

	err := envconfig.Process(SERVICENAME, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
