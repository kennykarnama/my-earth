package config

import "github.com/kelseyhightower/envconfig"

type Spec struct {
	AppName            string `envconfig:"APP_NAME" required:"true"`
	HTTPPort           string `envconfig:"HTTP_PORT" default:":8080"`
	DSN                string `envconfig:"DSN" required:"true"`
	MeteoSourceBaseURL string `envconfig:"METEO_SOURCE_BASE_URL" default:"https://www.meteosource.com/api/v1/free"`
	MeteoSourceAPIKey  string `envconfig:"METEO_SOURCE_API_KEY" required:"true"`
}

func Get() Spec {
	s := Spec{}
	envconfig.MustProcess("earth", &s)

	return s
}
