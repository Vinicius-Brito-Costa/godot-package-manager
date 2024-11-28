package util

import (
	"encoding/json"
	"fmt"
	"os"
)

type LoggingConfiguration struct {
	Level string
}
type Configuration struct {
	Logging LoggingConfiguration
}

func GetLoggingConfig() *LoggingConfiguration {
	var config, err = load()
	if err != nil {
		fmt.Println("error: ", err)
		return nil
	}

	return &config.Logging
}

func load() (*Configuration, error) {
	var file, err = os.Open("./config.json")
	if err != nil {
		fmt.Println("error: ", err)
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err = decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error: ", err)
		return nil, err
	}
	return &configuration, nil
}
