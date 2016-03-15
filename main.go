package main

import (
	_ "bytes"
	"fmt"
	_ "io/ioutil"
	"os"
	"os/exec"
	"regexp"
	_ "strings"

	"github.com/fatih/color"
	"github.com/golang/protobuf/proto"
)

const (
	mainFileContent = `
package main

import (
    "github.com/gin-gonic/gin"
    _ "net/http"
    "%s/test/handler"
    _ "%s/test/middleware"
    _ "%s/test/model"
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
{{ define "index.html" }}
<h1>Hello {{ .name }}</h1>
{{ end }}
`
	handlerNewFileContent = `
package handler

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

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
{{ define "%s" }}
{{ end }}
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
	r, _ := regexp.Compile(`(\w+)`)
	if r.FindString(apiName) != apiName {
		panic(fmt.Sprintf("\"%s\" is not allow", apiName))
	}
	newHandlerFilename := projectBaseDir + "/src/" + config.GetUser() + "/" + config.GetAppName() + "/handler/" + apiName + ".go"
	newModelFilename := projectBaseDir + "/src/" + config.GetUser() + "/" + config.GetAppName() + "/model/" + apiName + ".go"
	newTemplateFilename := projectBaseDir + "/src/" + config.GetUser() + "/" + config.GetAppName() + "/templates/" + apiName + ".html"
	WriteStringToFile(newHandlerFilename, fmt.Sprintf(handlerNewFileContent, apiName, apiName, apiName, apiName, apiName, apiName))
	WriteStringToFile(newModelFilename, modelNewFileContent)
	WriteStringToFile(newTemplateFilename, fmt.Sprintf(templateNewFileContent, apiName+".html"))
}

func middlerware(config *Config, middlerwareName string) {
	newMiddlerwareFilename := projectBaseDir + "/src/" + config.GetUser() + "/" + config.GetAppName() + "/middleware/" + middlerwareName + ".go"
	WriteStringToFile(newMiddlerwareFilename, middlewareNewFileContent)
}

func listDeps(config *Config) {
	deps := config.GetDeps()
	for _, dep := range deps {
		fmt.Println(dep.GetPkg())
	}
}

func addDeps(config *Config, pkg string) {
	deps := config.GetDeps()
	for _, dep := range deps {
		if dep.GetPkg() == pkg {
			panic(fmt.Sprintf("pkg: \"%s\" exists", pkg))
		}
	}
	fmt.Printf("Add pkg: \"%s\"\n", pkg)
	config.Deps = append(deps, &Dependency{
		Pkg: proto.String(pkg),
	})
	fmt.Println(config.Deps)
	WriteConfig(config, projectBaseDir+"/.beer.config")
}

func delDeps(config *Config, pkg string) {
	deps := config.GetDeps()
	for i, dep := range deps {
		if dep.GetPkg() == pkg {
			fmt.Printf("Delete pkg: \"%s\"\n", pkg)
			deps = append(deps[:i], deps[i+1:]...)
			config.Deps = deps
			WriteConfig(config, projectBaseDir+"/.beer.config")
			return
		}
	}
	panic(fmt.Sprintf("pkg: \"%s\" isn't in deps list", pkg))
}

func getDeps(config *Config) {
	deps := config.GetDeps()
	for _, dep := range deps {
		fmt.Printf("Installing pkg \"%s\"\n", dep.GetPkg())
		exec.Command("go", "get", dep.GetPkg()).Run()
	}
}

func buildApp(config *Config) {
	fmt.Println("Building handler")
	exec.Command("go", "build", fmt.Sprintf("%s/%s/handler", config.GetUser(), config.GetAppName())).Run()
	fmt.Println("Building middleware")
	exec.Command("go", "build", fmt.Sprintf("%s/%s/middleware", config.GetUser(), config.GetAppName())).Run()
	fmt.Println("Building model")
	exec.Command("go", "build", fmt.Sprintf("%s/%s/model", config.GetUser(), config.GetAppName())).Run()
	fmt.Println("Building main")
	exec.Command("go", "build", "-o", fmt.Sprintf("%s/src/%s/%s/main.out", projectBaseDir, config.GetUser(), config.GetAppName()), fmt.Sprintf("%s/%s/main", config.GetUser(), config.GetAppName())).Run()
}

func runApp(config *Config, mode string) {
	os.Setenv("GIN_MODE", mode)
	os.Chdir(projectBaseDir + "/src/" + config.GetUser() + "/" + config.GetAppName())
	cmd := exec.Command("./main.out")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
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

	os.Setenv("GOPATH", projectBaseDir)
	config := ReadConfig(projectBaseDir + ".beer.config")
	if command == "generate" {
		if len(Args) <= 2 {
			red := color.New(color.FgRed).SprintFunc()
			panic(red("error:  ") + "Usage: beer generate [something]")

		}
		apiName := Args[2]
		generate(config, apiName)
		return
	} else if command == "middleware" {
		if len(Args) <= 2 {
			red := color.New(color.FgRed).SprintFunc()

			yellow := color.New(color.FgYellow).SprintFunc()

			panic(red("error:  ") + yellow("Usage: beer middleware [something]"))

			return
		}
		middlerwareName := Args[2]
		middlerware(config, middlerwareName)
	} else if command == "run" {
		mode := "debug"
		if len(Args) > 2 {
			mode = Args[2]
		}
		fmt.Println("Run in mode:", mode)
		getDeps(config)
		buildApp(config)
		runApp(config, mode)
	} else if command == "deps" {
		if len(Args) <= 2 {
			red := color.New(color.FgRed).SprintFunc()
			panic(red("error:  ") + "Usage: beer deps [something]")
		}

		depsCmd := Args[2]
		if depsCmd == "list" {
			listDeps(config)
		} else if depsCmd == "add" {
			if len(Args) <= 3 {
				red := color.New(color.FgRed).SprintFunc()
				panic(red("error:  ") + "Usage: beer deps add [something]")
			}
			pkg := Args[3]
			addDeps(config, pkg)
		} else if depsCmd == "del" {
			if len(Args) <= 3 {
				red := color.New(color.FgRed).SprintFunc()
				panic(red("error:  ") + "Usage: beer deps add [something]")
			}
			pkg := Args[3]
			delDeps(config, pkg)
		} else if depsCmd == "get" {
			getDeps(config)
		}
	}
}
