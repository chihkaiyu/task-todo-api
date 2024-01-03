package config

type (
	Config struct {
		Env   string `env:"ENV" default:"local"`
		Port  string `env:"PORT" default:"8080"`
		Debug bool   `env:"DEBUG" default:"false"`
	}
)
