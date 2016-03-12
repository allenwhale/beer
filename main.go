package main

import (
	"fmt"
	_ "io/ioutil"
	"os"

	"github.com/fatih/color"
	_ "github.com/golang/protobuf/proto"
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
        c.String("Hello, go get beer")
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

		if len(Args) <= 2 {
			red := color.New(color.FgRed).SprintFunc()
			fmt.Println(red("error:  ") + "Usage: beer new [something]")
			return
		}
		config := NewConfig(Args[2])
		err := NewApp(config)
		if err != nil {
			fmt.Println(err)
		}
		return
	} else if command == "generate" {

		if len(Args) <= 2 {
			red := color.New(color.FgRed).SprintFunc()
			fmt.Println(red("error:  ") + "Usage: beer generate [something]")
			return
		}
		apiName := Args[2]
		fmt.Println(apiName)
		return
	}
}
