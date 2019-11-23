package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	// Sesitive data
	viper.SetEnvPrefix("ICS")
	viper.BindEnv("ENVIRONMENT")

	switch viper.Get("ENVIRONMENT") {
	case "DEVELOPMENT":
		log.SetLevel(log.DebugLevel)
	case "RELEASE":
		log.SetLevel(log.InfoLevel)
	}
	log.Infof("log configured to %s mode", log.GetLevel().String())

	// Static values
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Info("server configuration loaded")
}

func main() {
	// Running HTTP API
	addr := viper.GetString("addr")
	log.Infof("message server listening on %s", addr)
	if err := http.ListenAndServe(addr, EnableCORS(Router)); err != nil {
		log.Fatalln(err.Error())
	}
}
