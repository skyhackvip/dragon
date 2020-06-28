package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
)

type EnvConfig struct {
	Server struct {
		ListenAddress string `json:"listen_address"`
	} `json:"server"`
	Log struct {
		Access      string `json:"access"`
		Application string `json:"application"`
	} `json:"log"`
	Rabbitmq struct {
		Server string `json:"server"`
		Port   int    `json:"port"`
	} `json:"rabbitmq"`
	Es struct {
		Server string `json:"server"`
		Port   int    `json:"port"`
		Index  string `json:"index"`
		Type   string `json:"type"`
	} `json:"es"`
}

var GlobalEnv EnvConfig

func LoadConfig() {
	fmt.Println("loadconfig")
	var configFile = flag.String("c", "./config/env.json", "the path of the config file")
	flag.Parse()
	bytes, err := ioutil.ReadFile(*configFile)
	if err != nil {
		panic("load config error!")
	}
	err = json.Unmarshal(bytes, &GlobalEnv)
	if err != nil {
		panic(err)
	}
}
