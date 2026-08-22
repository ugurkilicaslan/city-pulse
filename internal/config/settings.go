package config

type Settings struct {
	Env string

	MongoURI      string
	MongoDatabase string

	GNewsAPIKey string
	NASAAPIKey  string
	AuthAPIKey  string

	ServerPort       int
	ServerAPIVersion string
}

func (s Settings) IsLocal() bool {
	return s.Env == "local"
}

func (s Settings) APIBasePath() string {
	version := s.ServerAPIVersion
	if version == "" {
		version = "1.0"
	}
	return "/api/" + version
}
