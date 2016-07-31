package main

import (
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"testing"
)

func TestEnvVarPopulation(t *testing.T) {
	config := Config{}

	os.Setenv("MY_TEST_VAR", "test")

	config.addEnvVar("MY_TEST_VAR")
	ctx, _ := newContext(&config)

	if ctx.Env["MY_TEST_VAR"] != "test" {
		t.Fatalf("test issue with '%s'\n", ctx.Env["MY_TEST_VAR"])
	}
}

func TestEnvVarPopulation_not_defined(t *testing.T) {
	config := Config{}
	config.addEnvVar("NOT_EXISTING_VAR")
	_, err := newContext(&config)

	if err == nil {
		t.Error("Error expected for undefined env var")
	} else if !strings.Contains(err.Error(), "NOT_EXISTING_VAR") {
		t.Error("Error message does not contain NOT_EXISTING_VAR, value=", err.Error())
	}
}

func TestTemplate(t *testing.T) {
	testTplPath1 := "./tests/tmp/mytemplate1.txt"
	testTplPath2 := "./tests/tmp/mytemplate2.txt"
	defer os.RemoveAll(filepath.Dir(testTplPath1))

	envVarValue := "foobar"

	os.Setenv("MY_TEST_VAR", envVarValue)

	config := Config{}
	config.addEnvVar("MY_TEST_VAR")

	currentUser, err := user.Current()
	if err != nil {
		t.Errorf("could not lookup current user, error: %v", err)
	}

	currentGroup, err := lookupGroupById(currentUser.Gid)

	config.addTemplate("./tests/test.tpl", testTplPath1, currentUser.Username, currentGroup.Name, "")
	config.addTemplate("./tests/test.tpl", testTplPath2, currentUser.Username, currentGroup.Name, "0500")
	ctx, err := newContext(&config)
	if err != nil {
		t.Errorf("newContext error: %v", err)
	}

	err = generateTemplates(config.Templates, ctx)
	if err != nil {
		t.Errorf("generateTemplates, error: %v", err)
	}

	assertFileExistsAndContains(testTplPath1, envVarValue, t)
	assertFileExistsAndContains(testTplPath2, envVarValue, t)
	assertFileMode(testTplPath2, 0500, t)
}

func TestMyMain(t *testing.T) {
	testTplPath1 := "./tests/dest/mytemplate.txt"
	defer os.RemoveAll(filepath.Dir(testTplPath1))

	envVarValue := "foobar"
	os.Setenv("MY_TEST_VAR", envVarValue)
	os.Args[1] = "-config=tests/config.json"

	main()

	assertFileExistsAndContains(testTplPath1, envVarValue, t)
	assertFileMode(testTplPath1, 0700, t)
}

func assertFileExistsAndContains(file, content string, t *testing.T) {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		t.Errorf("file %s does not exist", file)
	} else {
		b, err := ioutil.ReadFile(file)
		if err != nil {
			t.Error("file %s could not be read, error: %v", file, err)
		}
		content := string(b)
		if !strings.Contains(content, content) {
			t.Error("File %s does not contain ", file, content)
		}
	}
}

func assertFileMode(file string, fileMode os.FileMode, t *testing.T) {
	fileStat, err := os.Stat(file)
	if err != nil {
		t.Errorf("File %s stat error %v", file, err)
	}

	if fileStat.Mode() != fileMode {
		t.Errorf("unexpected mode %v for file %s", fileStat.Mode().Perm(), file)
	}
}
