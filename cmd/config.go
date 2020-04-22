package main

import (
	"bytes"
	"fmt"
	"github.com/vidmed/request"
	"time"

	"github.com/BurntSushi/toml"
)

var configInstance *TomlConfig

// TomlConfig represents a config
type TomlConfig struct {
	Main Main
}

// Main represent a main section of the TomlConfig
type Main struct {
	LogLevel   uint8
	ListenAddr string
	ListenPort uint
}

// GetConfig returns application config
func GetConfig() *TomlConfig {
	return configInstance
}

// NewConfig creates new application config with given .toml file
func NewConfig(file string) (*TomlConfig, error) {
	configInstance = &TomlConfig{}

	if _, err := toml.DecodeFile(file, configInstance); err != nil {
		return nil, err
	}
	dump(configInstance)

	// check required fields
	// Main
	if configInstance.Main.ListenAddr == "" {
		request.GetLogger().Fatalln("Main.ListenAddr must be specified. Check your Config file")
	}
	if configInstance.Main.ListenPort == 0 {
		request.GetLogger().Fatalln("Main.ListenPort must be specified. Check your Config file")
	}
	return configInstance, nil
}

func dump(cfg *TomlConfig) {
	var buffer bytes.Buffer
	e := toml.NewEncoder(&buffer)
	err := e.Encode(cfg)
	if err != nil {
		request.GetLogger().Fatal(err)
	}

	fmt.Println(
		time.Now().UTC(),
		"\n---------------------Sevice started with config:\n",
		buffer.String(),
		"\n---------------------")
}
