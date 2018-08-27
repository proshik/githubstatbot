package config

import "github.com/kelseyhightower/envconfig"

const (
	// SERVICENAME contains a service name prefix which used in ENV variables
	SERVICENAME = "GITHUBSTATBOT"
)

type Config struct {
	Mode               string `default:"prod"`
	Port               string `envconfig:"PORT",default:"8080"`
	TlsDir             string `default:"./"`
	LogDir             string `default:"./log/"`
	StaticFilesDir     string `default:"./static"`
	DbPath             string `default:"database.db"`
	GitHubClientId     string `required:"true"`
	GitHubClientSecret string `required:"true"`
	TelegramToken      string `required:"true"`
	AuthBasicUsername  string `default:"user"`
	AuthBasicPassword  string `default:"password"`
	DbUrl              string `envconfig:"JDBC_DATABASE_URL",default:"postgres://postgres:password@localhost:5432/githubstatbot?sslmode=disable"`
}

func Load() (*Config, error) {

	cfg := new(Config)

	err := envconfig.Process(SERVICENAME, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
