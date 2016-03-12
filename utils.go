package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/golang/protobuf/proto"
)

func RecoverPanic() {
	err := recover()
	if err != nil {
		fmt.Println(err)
	}
}

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

func FindBaseDir() string {
	base := "."
	for ; ; base += "/.." {
		files, err := ioutil.ReadDir(base)
		Check(err)
		for _, file := range files {
			if file.Name() == ".beer.config" {
				base, err = filepath.Abs(base)
				Check(err)
				return base
			}
		}
		if p, _ := filepath.Abs(base); p == "/" {
			Check(errors.New(".beer.config Not found"))
		}
	}
	Check(errors.New("Something error"))
	return ""
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
		Deps: []*Dependency{
			&Dependency{
				Pkg: proto.String("github.com/gin-gonic/gin"),
			},
		},
	}
	return config
}

func ReadConfig(configPath string) *Config {
	in, err := ioutil.ReadFile(configPath)
	Check(err)
	config := &Config{}
	err = proto.Unmarshal(in, config)
	Check(err)
	return config
}
func WriteConfig(config *Config, filePath string) {
	out, err := proto.Marshal(config)
	Check(err)
	err = ioutil.WriteFile(filePath, out, 0644)
	Check(err)
}
func WriteStringToFile(filePath string, content string) {
	err := ioutil.WriteFile(filePath, []byte(content), 0644)
	Check(err)
}
