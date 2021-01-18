package main

import (
    "encoding/json"
    "fmt"
    "github.com/go-resty/resty/v2"
    "github.com/shumybest/ragnaros2/config"
    "github.com/urfave/cli/v2"
    "io"
    "io/ioutil"
    "os"
    "path"
    "path/filepath"
    "sort"
    "strings"
    "text/template"
)

var Commands = []*cli.Command{
    {
        Name:    "version",
        Aliases: []string{"ver"},
        Usage:   "tool version",
        Action: func(c *cli.Context) error {
            fmt.Println(config.GetVersion())
            return nil
        },
    },
    {
        Name:    "download",
        Aliases: []string{"down"},
        Action: func(c *cli.Context) error {
            filename, _ := os.Executable()
            binDir, _ := filepath.Abs(filepath.Dir(filename))
            downloadTemplates(binDir, c)
            return nil
        },
        Flags: []cli.Flag{
            &cli.BoolFlag{
                Name:    "force",
                Aliases: []string{"f"},
                Usage:   "force to download and overwrite template files",
            },
        },
    },
    {
        Name:    "generate",
        Aliases: []string{"gen"},
        Usage:   "generate a micro service project with initial codes",
        Action:  generateProject,
        Flags: []cli.Flag{
            &cli.StringFlag{
                Name:    "output",
                Aliases: []string{"o"},
                Value:   "./",
                Usage:   "output directory of the generated project",
            },
            &cli.StringFlag{
                Name:     "conf",
                Aliases:  []string{"c"},
                Value:    "",
                Required: true,
                Usage: `load configuration data from FILE, FILE Content refer this example:
    {
		"App": {
			"ProjectName": "exampleProject",
			"ControllerName": "exampleController"
		},
        "K8s": {
	        "Spring": {
	            "Profiles": "prod",
	            "CloudConfigUri": "http://admin:${registry.password}@example.cluster.local:8761/config",
	            "DataSourceUrl": "jdbc:mysql://example.mysql.rds.aliyuncs.com:3306/example?useUnicode=true&characterEncoding=utf8&useSSL=false&useLegacyDatetimeCode=false&serverTimezone=Asia/Shanghai"
	        },
	        "Server": {
	            "Port": 8999
	        },
	        "Eureka": {
	            "ServiceUrl": "http://admin:${registry.password}@example.cluster.local:8761/eureka/"
	        }
        }
    }
`,
            },
        },
    },
}

var downloadSource = "https://raw.githubusercontent.com/shumybest/ragnaros2/master/cmd/ragnaros/"
var templateFiles = map[string]string{
    "tpls/Makefile.tpl":                              "Makefile",
    "tpls/Dockerfile.tpl":                            "Dockerfile",
    "tpls/go.mod.tpl":                                "go.mod",
    "tpls/k8s-deployment.yml.tpl":                    "k8s-deployment.yml",
    "tpls/main.go.tpl":                               "main.go",
    "tpls/app/controller.go.tpl":                     "app/controller.go",
    "tpls/resources/config/bootstrap.yml.tpl":        "resources/config/bootstrap.yml",
    "tpls/resources/config/bootstrap-prod.yml.tpl":   "resources/config/bootstrap-prod.yml",
    "tpls/resources/config/application.yml.tpl":      "resources/config/application.yml",
    "tpls/resources/config/application-dev.yml.tpl": "resources/config/application-dev.yml",
    "tpls/resources/config/application-prod.yml.tpl": "resources/config/application-prod.yml",
}

func downloadTemplates(binDir string, c *cli.Context) {
    for key, _ := range templateFiles {
        if _, err := os.Stat(path.Join(binDir, key)); os.IsNotExist(err) || c.Bool("force") {
            fmt.Println("Downloading template file : " + key + " into " + binDir)

            client := resty.New()
            client.SetHostURL(downloadSource)
            if resp, err := client.R().Get(key); err == nil {
                absOutputFileName := path.Join(binDir, key)
                _ = os.MkdirAll(path.Dir(absOutputFileName), os.ModePerm)
                outputFile, _ := os.Create(absOutputFileName)
                _, _ = io.WriteString(outputFile, resp.String())
                outputFile.Close()
            } else {
                fmt.Println(err)
                continue
            }
        }
    }
}

func generateFiles(binDir string, data interface{}, outputDirectory string) {
    for key, value := range templateFiles {
        tpl, err := template.New(path.Base(key)).
            Funcs(map[string]interface{}{
                "Export": strings.Title}).ParseFiles(path.Join(binDir, key))
        if err != nil {
            fmt.Println(err.Error())
            continue
        }

        absOutputFileName := path.Join(outputDirectory, value)
        _ = os.MkdirAll(path.Dir(absOutputFileName), os.ModePerm)
        outputFile, err := os.Create(absOutputFileName)
        fmt.Println("Generating file : " + absOutputFileName)

        err = tpl.Execute(outputFile, data)
        if err != nil {
            panic(err)
        }
    }
}

func generateProject(c *cli.Context) error {
    filename, _ := os.Executable()
    binDir, _ := filepath.Abs(filepath.Dir(filename))

    downloadTemplates(binDir, c)

    dataFileName := c.String("conf")
    outputDirectory := c.String("output")

    dataFile, err := os.Open(dataFileName)
    if err != nil {
        panic(err)
    }
    defer dataFile.Close()

    byteValue, err := ioutil.ReadAll(dataFile)
    var data interface{}
    err = json.Unmarshal(byteValue, &data)
    if err != nil {
        panic(err)
    }

    generateFiles(binDir, data, outputDirectory)
    return nil
}

func main() {
    cliApp := cli.NewApp()
    cliApp.Name = "ragnaros framework helper"
    cliApp.Usage = "[-h]"
    cliApp.Commands = Commands
    sort.Sort(cli.CommandsByName(cliApp.Commands))

    err := cliApp.Run(os.Args)
    if err != nil {
        fmt.Println(err)
    }
}
