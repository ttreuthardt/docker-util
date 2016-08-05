package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

type Template struct {
	TemplatePath string `json:"templatePath"`
	DestPath     string `json:"destPath"`
	FileMode     string `json:"fileMode"`
	DirMode      string `json:"dirMode"`
	Owner        string `json:"owner"`
	Group        string `json:"group"`
}

type Config struct {
	Envvars   []string   `json:"envvars"`
	Templates []Template `json:"templates"`
}

func (c *Config) addEnvVar(envVarName string) {
	c.Envvars = append(c.Envvars, envVarName)
}

func (c *Config) clearTemplates() {
	c.Templates = make([]Template, 0)
}

func (c *Config) addTemplate(templatePath, destPath, owner, group, fileMode string) error {

	for _, t := range c.Templates {
		if t.TemplatePath == templatePath {
			errors.New(fmt.Sprintf("The TemplatePath '%s' exists already", templatePath))
		}
	}

	c.Templates = append(c.Templates, Template{
		TemplatePath: templatePath,
		DestPath:     destPath,
		Owner:        owner,
		Group:        group,
		FileMode:     fileMode,
	})

	return nil
}

func readConfig(configFilePath string) (*Config, error) {
	file, e := ioutil.ReadFile(configFilePath)
	if e != nil {
		return nil, errors.New(fmt.Sprintf("Error reading config file '%s', error: %v\n", configFilePath, e))
	}
	var config Config
	json.Unmarshal(file, &config)

	return &config, nil
}
