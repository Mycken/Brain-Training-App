package main

import (
	"BrainTApp/internal/bta/apiserver"
	"flag"
	"github.com/BurntSushi/toml"
	"log"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "path to config file")
}

func main() {
	flag.Parse()

	config := apiserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	//s := apiserver.New(config)
	if err := apiserver.Start(config); err != nil {

	}

}
