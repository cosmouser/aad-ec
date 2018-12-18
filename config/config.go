package config

import (
	"encoding/json"
	"flag"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

// Config contains domain specific credentials
type Config struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Resource     string `json:"resource"`
	GrantType    string `json:"grant_type"`
	Tenant       string `json:"tenant"`
	APIVersion   string `json:"api-version"`
	ExternalURL  string `json:"externalURL"`
}

// C holds data from the parsed config file for use throughout the server
var C Config

// Port is the port number that the server listens on
var Port string

var configPath = flag.String("config", "./config.json", "path to config.json")
var tmp = flag.Int("port", 8080, "port to listen on")

func init() {
	flag.Parse()
	Port = ":" + strconv.Itoa(*tmp)
	file, err := os.Open(*configPath)
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&C); err != nil {
		log.Fatal(err)
	}
	file.Close()
}
