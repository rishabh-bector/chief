package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

var globalConfig ChiefConfig

// Security clearances
const MASTER_CLEARANCE = "master"
const NORMAL_CLEARANCE = "normal"

type ChiefConfig struct {
	Access map[string]User `json:"access"`
}

// Setup is run on installation, to create the basic config
func Setup() error {
	globalConfig = ChiefConfig{
		Access: make(map[string]User),
	}

	err := WriteToDisk()
	if err != nil {
		return err
	}

	return nil
}

type User struct {
	PassHash  string `json:"hash"`
	Clearance string `json:"clearance"`
}

func Global() *ChiefConfig {
	return &globalConfig
}

func Ensure() error {
	if globalConfig.Access != nil {
		return nil
	}

	err := LoadFromDisk()
	if err != nil {
		return errors.New("unable to load chief config, have you run 'chief setup' yet?")
	}
	return nil
}

func LoadFromDisk() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	b, err := ioutil.ReadFile(fmt.Sprintf("%s/.chief/config.json", home))
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, &globalConfig)
	if err != nil {
		return err
	}

	return nil
}

func WriteToDisk() error {
	b, err := json.MarshalIndent(globalConfig, "", "    ")
	if err != nil {
		return err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	f, err := os.Create(fmt.Sprintf("%s/.chief/config.json", home))
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(b)
	if err != nil {
		return err
	}

	return nil
}
