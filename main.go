package main

import (
	"bytes"
	"fmt"
	_ "io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
)

const (
	mainFileContent = `
package main

import (
    "github.com/gin-gonic/gin"
    _ "net/http"
    "%s/test/handler"
    "%s/test/middleware"
    "%s/test/model"
)

func main() {
    router := gin.New()
    router.Use(gin.Logger())
    router.Use(gin.Recovery())

    router.LoadHTMLGlob("template/**/*")
    router.Static("/static", "./static")
    router.GET("/", handler.IndexGET)
    router.Run(":8080")
    return
}
`
	handlerMainFileContent = `
package handler
`
	handlerIndexFileContent = `
package handler

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

func IndexGET(c *gin.Context) {
    c.HTML(http.StatusOK, "index.html", gin.H{"name": "index"})
}
`
	modelMainFileContent = `
package model
`
	middlerwareMainFileContent = `
package middlerware
`
	templateIndexFileContent = `
<h1>Hello {{ .name }}</h1>
`
	handlerNewFileContent = `
package handler


func %sGET(c *gin.Context) {
}
func %sPOST(c *gin.Context) {
}
func %sPUT(c *gin.Context) {
}
func %sDELETE(c *gin.Context) {
}
func %sPATCH(c *gin.Context) {
}
func %sOPTION(c *gin.Context) {
}
`
	modelNewFileContent = `
package model
`
	middlewareNewFileContent = `
package middleware
`
	templateNewFileContent = `
`
)

var projectBaseDir string

func NewApp(config *Config) {

	baseDir := config.GetAppName() + "/src/" + config.GetUser() + "/" + config.GetAppName()
	mainFolder := baseDir + "/main/"
	handlerFolder := baseDir + "/handler/"
	middlerwareFolder := baseDir + "/middleware/"
	modelFolder := baseDir + "/model/"
	templateFolder := baseDir + "/template/"
	staticFolder := baseDir + "/static/"
	folders := []string{mainFolder, handlerFolder, modelFolder, templateFolder, middlerwareFolder, staticFolder}
	for _, folder := range folders {
		err := os.MkdirAll(folder, 0755)
		Check(err)
	}
	WriteConfig(config, config.GetAppName()+"/.beer.config")
	WriteStringToFile(baseDir+"/main/main.go", fmt.Sprintf(mainFileContent, config.GetUser(), config.GetUser(), config.GetUser()))
	WriteStringToFile(baseDir+"/handler/main.go", handlerMainFileContent)
	WriteStringToFile(baseDir+"/handler/index.go", handlerIndexFileContent)
	WriteStringToFile(baseDir+"/middleware/main.go", middlerwareMainFileContent)
	WriteStringToFile(baseDir+"/model/main.go", modelMainFileContent)
	WriteStringToFile(baseDir+"/template/index.html", templateIndexFileContent)
	fmt.Println("Run:")
	fmt.Printf("cd %s\n", config.GetAppName())
	fmt.Println("export GOPATH=`pwd`")
}

func generate(config *Config, apiName string) {
	apiUrl := strings.Split(apiName, "/")
	apiName = ""
	for _, url := range apiUrl {
		apiName += strings.ToUpper(url[0:1]) + url[1:]
	}
	newHandlerFilename := projectBaseDir + "/src/" + config.GetUser() + "/" + config.GetAppName() + "/handler/" + apiName + ".go"
	newModelFilename := projectBaseDir + "/src/" + config.GetUser() + "/" + config.GetAppName() + "/model/" + apiName + ".go"
	newTemplateDir := projectBaseDir + "/src/" + config.GetUser() + "/" + config.GetAppName() + "/template/" + apiName + "/"
	newTemplateFilename := newTemplateDir + "index.html"
	err := os.MkdirAll(newTemplateDir, 0755)
	Check(err)
	WriteStringToFile(newHandlerFilename, fmt.Sprintf(handlerNewFileContent, apiName, apiName, apiName, apiName, apiName, apiName))
	WriteStringToFile(newModelFilename, modelNewFileContent)
	WriteStringToFile(newTemplateFilename, templateNewFileContent)
}

func middlerware(config *Config, middlerwareName string) {
	newMiddlerwareFilename := projectBaseDir + "/src/" + config.GetUser() + "/" + config.GetAppName() + "/middleware/" + middlerwareName + ".go"
	WriteStringToFile(newMiddlerwareFilename, middlewareNewFileContent)
}

func buildApp(config *Config) {
	fmt.Println(os.Getenv("GOPATH"))
	var out bytes.Buffer
	cmd := exec.Command("go", "build", fmt.Sprintf("%s/%s/handler", config.GetUser(), config.GetAppName()))
	cmd.Run()
	exec.Command("go", "build", fmt.Sprintf("%s/%s/middleware", config.GetUser(), config.GetAppName())).Run()
	exec.Command("go", "build", fmt.Sprintf("%s/%s/model", config.GetUser(), config.GetAppName())).Run()
	cmd = exec.Command("go", "-o", "main.out", "build", fmt.Sprintf("%s/%s/main", config.GetUser(), config.GetAppName()))
	cmd.Stdout = &out
	err := cmd.Run()
	fmt.Println(err)
	fmt.Printf("in all caps: %q\n", out.String())
}

func main() {
	defer RecoverPanic()
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
			yellow := color.New(color.FgYellow).SprintFunc()
			fmt.Println(red("error:  ") + yellow("Usage: beer new [something]"))
			return
		}

		if _, err := os.Stat(Args[2]); err == nil {
			red := color.New(color.FgRed).SprintFunc()
			yellow := color.New(color.FgYellow).SprintFunc()
			fmt.Printf(red("error: ") + yellow(Args[2]+" already exists! Override? (y/n) [n]: "))
			var opt string
			fmt.Scanf("%s", &opt)
			if opt != "y" {
				return
			}
			exec.Command("rm", "-rf", Args[2]).Run()
		}

		config := NewConfig(Args[2])
		NewApp(config)
		green := color.New(color.FgGreen).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()
		fmt.Println(green("Your beer is ready at ") + yellow(config.GetAppName()+"/src/"+config.GetUser()))
		return
	}
	projectBaseDir = FindBaseDir() + "/"

	if command == "generate" {
		if len(Args) <= 2 {
			red := color.New(color.FgRed).SprintFunc()
			yellow := color.New(color.FgYellow).SprintFunc()
			fmt.Println(red("error:  ") + yellow("Usage: beer generate [something]"))
			return
		}
		apiName := Args[2]
		config := ReadConfig(projectBaseDir + ".beer.config")
		generate(config, apiName)
		return
	} else if command == "middleware" {
		if len(Args) <= 2 {
			red := color.New(color.FgRed).SprintFunc()
			yellow := color.New(color.FgYellow).SprintFunc()
			fmt.Println(red("error:  ") + yellow("Usage: beer middleware [something]"))
			return
		}
		middlerwareName := Args[2]
		config := ReadConfig(projectBaseDir + ".beer.config")
		middlerware(config, middlerwareName)
	} else if command == "run" {
		mode := "debug"
		if len(Args) > 2 {
			mode = Args[3]
		}
		fmt.Println("Run in mode: ", mode)
		config := ReadConfig(projectBaseDir + ".beer.config")
		buildApp(config)
	}
}
