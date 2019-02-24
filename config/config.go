package config

import (
	"github.com/golang/protobuf/proto"
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func ReadAndMaybeCreateConfig() *Configuration {
	configDir, _ := homedir.Expand("~/.config/godict")
	configPath := filepath.Join(configDir, "config.textproto")

	// Attempt to create config dir, in case it's missing.
	os.MkdirAll(configDir, os.ModePerm)

	config := &Configuration{}
	data, err := ioutil.ReadFile(configPath)

	if err == nil {
		err = proto.UnmarshalText(string(data), config)
		if err == nil {
			log.Printf("Read config file from %s.\n", configPath)
		} else {
			log.Fatalf("Error reading config file from %s: %s.\n", configPath, err)
		}
	} else {
		config.StorePath = proto.String("~/.godict-store")
		config.MaxNonpersistedChanges = proto.Int32(4)

		// Regardless of the error type, we'll attempt to create the config.
		log.Printf("Couldn't read config file from %s. Using defaults:\n",
			configPath)
		log.Println(config)

		err = ioutil.WriteFile(configPath, []byte(proto.MarshalTextString(config)), 0644)
		if err != nil {
			log.Printf("Couldn't write config to %s: %s\n", configPath, err)
		}
	}

	return config
}
