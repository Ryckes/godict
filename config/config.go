
package config

import (
	"github.com/golang/protobuf/proto"
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
	"log"
)

func ReadConfig() *Configuration {
	configPath, _ := homedir.Expand("~/.config/godict/config.textproto")
	config := &Configuration{}
	data, err := ioutil.ReadFile(configPath)

	if err != nil {
		log.Fatalf("Couldn't read config file: %s.\n", err)
	}

	proto.UnmarshalText(string(data), config)
	return config
}
