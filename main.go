package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"text/template"
)

type Context struct {
	Env map[string]string
}

func main() {
	configFilePath := getConfigFilePath()

	config, err := readConfig(configFilePath)
	checkExitError(err)

	ctx, err := newContext(config)
	checkExitError(err)

	err = generateTemplates(config.Templates, ctx)
	checkExitError(err)
}

func checkExitError(err error) {
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}
}

// creates a new context and initializes it with the Env vars defined in the configuration
func newContext(config *Config) (*Context, error) {
	ctx := &Context{Env: make(map[string]string)}

	for _, envVarName := range config.Envvars {
		value := os.Getenv(envVarName)
		if value != "" {
			//log.Printf("Env var %s found with value '%s'", envVarName, value)
			ctx.Env[envVarName] = value
		} else {
			return nil, errors.New(fmt.Sprintf("Env var %s not defined!", envVarName))
		}
	}

	return ctx, nil
}

// reads the flag and returns config file path
func getConfigFilePath() string {
	var configFile string
	flag.StringVar(&configFile, "config", "./config.json", "JSON config file path")
	flag.Parse()

	log.Printf("Using config file %s", configFile)

	return configFile
}

func tplFuncTest() string {
	return "testoutput"
}

func generateTemplates(templates []Template, ctx *Context) error {
	for _, template := range templates {
		err := handleTemplate(template, ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func handleTemplate(template Template, ctx *Context) error {
	var err error
	dest := os.Stdout //TODO: add proper way to print templates to stdout
	if template.DestPath != "" {
		createFolderStructure(template.DestPath)
		dest, err = os.Create(template.DestPath)
		if err != nil {
			return fmt.Errorf("could not create destination file %s, error: %v", template.DestPath, err)
		}
		defer dest.Close()
	}

	if err = writeTemplate(template.TemplatePath, dest, ctx); err != nil {
		return err
	}

	if err = handlePermission(dest, template); err != nil {
		return err
	}

	if err = handleOwnerAndGroup(dest, template); err != nil {
		return err
	}

	return nil
}

// writes the given template to the target destination
func writeTemplate(templatePath string, writer io.Writer, ctx *Context) error {
	tmpl, err := template.New(filepath.Base(templatePath)).Funcs(template.FuncMap{
		"test": tplFuncTest,
	}).ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("Failed to prepare template %s, error: %v", templatePath, err)
	}

	err = tmpl.Execute(writer, &ctx)
	if err != nil {
		return fmt.Errorf("Failed to generate template %s, error: %v", templatePath, err)
	}
	return nil
}

func handlePermission(file *os.File, template Template) error {
	if template.FileMode != "" {

		if !isValidFileMode(template.FileMode) {
			return fmt.Errorf("Invalid file mode %s for template %s", template.FileMode, template.DestPath)
		}

		if value, err := strconv.ParseUint(template.FileMode, 0, 32); err == nil {
			fileMode := os.FileMode(value)
			err := os.Chmod(file.Name(), fileMode)
			if err != nil {
				return fmt.Errorf("Chmod with mode %v failed for template destPath %s, error: %v", fileMode,
					file.Name(), err)
			}
		}
	}
	return nil
}

func handleOwnerAndGroup(file *os.File, template Template) error {
	uid := os.Geteuid()
	gid := os.Getegid()

	if template.Owner != "" {
		if isNumber(template.Owner) {
			uid, _ = strconv.Atoi(template.Owner)
		} else {
			owner, err := user.Lookup(template.Owner)
			if err == nil {
				uid, _ = strconv.Atoi(owner.Uid)
			} else {
				return fmt.Errorf("User %s not found chown to curent user, error: %v", template.Owner, err)
			}
		}
	}

	if template.Group != "" {
		if isNumber(template.Group) {
			gid, _ = strconv.Atoi(template.Group)
		} else {
			group, err := LookupGroupByName(template.Group)
			if err == nil {
				gid, _ = strconv.Atoi(group.Gid)
			} else {
				return fmt.Errorf("Group %s not found chown to curent users primary group, error: %v", template.Group, err)
			}
		}
	}

	err := file.Chown(uid, gid)
	if err != nil {
		return fmt.Errorf("Chown with uid %d and gid %d failed for template destPath %s, error: %v", uid, gid, template.DestPath, err)
	}

	return nil
}

// creates the parent folder it it does not exist
func createFolderStructure(templatePath string) error {
	path := filepath.Dir(templatePath)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755) //TODO: use dir mode if available
		if err != nil {
			fmt.Errorf("Could not create dir structure for template destPath %s, error: %v", templatePath, err)
		}
	}
	return nil
}

func isNumber(id string) bool {
	_, err := strconv.Atoi(id)
	return err == nil
}

func isValidFileMode(fileMode string) bool {
	match, _ := regexp.MatchString("^[0-7]{4}$", fileMode)
	return match
}
