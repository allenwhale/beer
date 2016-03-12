package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/golang/protobuf/proto"
)

func FindBaseDir() (string, error) {
	base := "."
	for ; ; base += "/.." {
		files, err := ioutil.ReadDir(base)
		if err != nil {
			return "", err
		}
		for _, file := range files {
			if file.Name() == ".beer.config" {
				return filepath.Abs(base)
			}
		}
		if p, _ := filepath.Abs(base); p == "/" {
			return "", errors.New(".beer.config Not found")
		}
	}
	return "", nil
}

func NewConfig(appName string) *Config {
	var user string
	fmt.Printf("User: ")
	fmt.Scanf("%s", &user)
	// fmt.Printf("AppName: ")
	// fmt.Scanf("%s", &appName)
	config := &Config{
		User:    proto.String(user),
		AppName: proto.String(appName),
	}
	return config
}

func ReadConfig(configPath string) (*Config, error) {
	in, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	if err := proto.Unmarshal(in, config); err != nil {
		return nil, err
	}
	return config, nil
}
func WriteConfig(config *Config, filePath string) error {
	out, err := proto.Marshal(config)
	if err != nil {
		return nil
	}
	if err := ioutil.WriteFile(filePath, out, 0644); err != nil {
		return err
	}
	return nil
}
