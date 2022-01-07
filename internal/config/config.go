package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
)

type Config struct {
	Broker   string `json:"broker"`
	Username string `json:"username"`
	Password string `json:"password"`
	Topics   topics `json:"topics"`
}

func GetConfig() Config {
	cfg := Config{}

	configFilename := flag.String("c", "config.json", "Path to the json configuration file")

	flag.Parse()

	handleFile(*configFilename, &cfg)

	return cfg
}

func handleFile(filename string, cfg *Config) {
	jsonFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err)
	}

	err = json.Unmarshal([]byte(jsonFile), &cfg)
	if err != nil {
		log.Println(err)
	}
}
