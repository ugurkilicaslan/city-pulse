package config

import (
	"flag"
	"log"

	"github.com/spf13/viper"
)

var Config Settings
var DevMode bool

func FlagConfig() {
	dev := flag.Bool("dev", false, "local config kullan (config/local-config.json)")
	flag.Parse()
	DevMode = *dev
	if *dev {
		viper.SetConfigName("local-config")
	} else {
		viper.SetConfigName("prod-config")
	}
}

func ViperRead() {
	viper.SetConfigType("json")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("config okunamadı: %v", err)
	}

	Config = Settings{
		Env:              viper.GetString("env"),
		MongoURI:         viper.GetString("mongo.uri"),
		MongoDatabase:    viper.GetString("mongo.database"),
		GNewsAPIKey:      viper.GetString("gnews.apiKey"),
		NASAAPIKey:       viper.GetString("nasa.apiKey"),
		AuthAPIKey:       viper.GetString("auth.apiKey"),
		JWTSecret:        viper.GetString("auth.jwtSecret"),
		ServerPort:       viper.GetInt("server.port"),
		ServerAPIVersion: viper.GetString("server.apiVersion"),
	}
}
