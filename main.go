package main

import (
	"fmt"
	_ "github.com/golang/protobuf/proto"
	_ "io/ioutil"
	"os"
)

const (
	mainFileContent = `
package main

import (
    "github.com/gin-gonic/gin"
)

func main() {
    router := gin.New()
    router.Use(gin.Logger())
    router.Use(gin.Recovery())
    router.GET("/", func(c *gin.Context) {
        c.String("index")
    })
    router.Run(":8080")
    return
}
`
)

func NewApp(config *Config) error {
	mainFolder := config.GetAppName() + "/src/" + config.GetUser() + "/" + config.GetAppName() + "/main/"
	controllerFolder := config.GetAppName() + "/src/" + config.GetUser() + "/" + config.GetAppName() + "/controller/"
	middlerwareFolder := config.GetAppName() + "/src/" + config.GetUser() + "/" + config.GetAppName() + "/middlerware/"
	modelFolder := config.GetAppName() + "/src/" + config.GetUser() + "/" + config.GetAppName() + "/model/"
	templateFolder := config.GetAppName() + "/src/" + config.GetUser() + "/" + config.GetAppName() + "/template/"
	folders := []string{mainFolder, controllerFolder, modelFolder, templateFolder, middlerwareFolder}
	for _, folder := range folders {
		err := os.MkdirAll(folder, 0755)
		if err != nil {
			return err
		}
	}
	err := WriteConfig(config, config.GetAppName()+"/.beer.config")
	if err != nil {
		return err
	}
	return nil
}

func main() {
	fmt.Printf(mainFileContent)
	Args := os.Args
	if len(Args) <= 1 {
		Usage()
		return
	}
	command := Args[1]
	if command == "help" {
		Usage()
		return
	} else if command == "new" {
		config := NewConfig()
		err := NewApp(config)
		if err != nil {
			fmt.Println(err)
		}
		return
	}
	config, err := ReadConfig(".beer.config")
	if err != nil {
		config := NewConfig()
		WriteConfig(config, ".beer.config")
		fmt.Println(config.GetUser(), config.GetAppName())
	} else {
		fmt.Println(config.GetUser(), config.GetAppName())
	}

}
