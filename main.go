package main

import (
	"io"
	"os"
	"log"
	"path/filepath"
	"text/template"
)

func main() {
	print("test")
}

type Config struct {
	Templates []string
	Envvars []string
}

func (c *Config) addEnvVar(envVarName string) {
	c.Envvars = append(c.Envvars, envVarName)
}

func NewContext() *Context {
	return &Context{Env: make(map[string]string)}
}


type Context struct {
	Env map[string]string
}

// reads all flags and returns a context object
func parsFlags() {

}

// Reads the defined env vars from the environment and stores it in the context
func populateContext(ctx *Context, config Config) {
	for _, envVarName := range config.Envvars {
		value := os.Getenv(envVarName)
		ctx.Env[envVarName] = value
		log.Printf("%s=%s", envVarName, value)
	}
}

func tplFuncTest() string {
	return "testoutput"
}

func handleTemplate(templatePath string, destPath string, ctx *Context) {
	dest := os.Stdout
	if destPath != "" {
		dest, err := os.Create(destPath)
		if err != nil {
			log.Fatalf("could not create destination file %s", err)
		}
		createFolderStructure(destPath)
		defer dest.Close()
	}

	writeTemplate(templatePath, dest, ctx)

}

// writes the given template to the target destination
func writeTemplate(templatePath string, writer io.Writer, ctx *Context) bool {
	tmpl, err := template.New(filepath.Base(templatePath)).Funcs(template.FuncMap{
		"test": tplFuncTest,
	}).ParseFiles(templatePath)
	if err != nil {
		log.Fatal("error while parsing template, %v", err)
		return false;
	}

	err = tmpl.Execute(writer, &ctx)
	if err != nil {
		log.Fatal("error executing parsing template, %v", err)
		return false;
	}

	return true
}

func createFolderStructure(templatePath string) bool {
	path := filepath.Dir(templatePath)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0777)
		if err != nil {
			return false
		}
	}
	return true
}