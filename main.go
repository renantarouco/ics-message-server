package main

import (
	"net/http"

	"github.com/renantarouco/ics-message-server/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	// Sesitive data
	viper.SetEnvPrefix("ICS")
	viper.BindEnv("ENVIRONMENT")
	viper.BindEnv("JWT_KEY")

	switch viper.GetString("ENVIRONMENT") {
	case "DEVELOPMENT":
		log.SetLevel(log.DebugLevel)
	case "RELEASE":
		log.SetLevel(log.InfoLevel)
	}
	log.Infof("log configured to %s mode", log.GetLevel().String())

	log.Infof("ICS_ENVIRONMENT = %s", viper.GetString("ENVIRONMENT"))
	log.Debugf("ICS_JWT_KEY     = %s", viper.GetString("JWT_KEY"))

	// Static values
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Info("server configuration loaded")

	server.Init()
}

func main() {
	id := pflag.UintP("id", "i", 1, "Node ID")
	port := pflag.UintP("port", "p", 7000, "Client access port")
	cluster := pflag.StringSliceP("cluster", "c", []string{"http://127.0.0.1:7000"}, "Cluster peers")
	pflag.Parse()
	log.Infof("NODE_ID =     %d", *id)
	log.Infof("CLIENT_PORT = %d", *port)
	log.Infof("CLUSTER =     %v", *cluster)
	// Running HTTP API
	addr := viper.GetString("addr")
	log.Infof("message server listening on %s", addr)
	if err := http.ListenAndServe(addr, EnableCORS(Router)); err != nil {
		log.Fatalln(err.Error())
	}
}
