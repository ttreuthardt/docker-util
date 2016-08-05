package main

import (
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
)

var (
	config      = Config{}
	envVarValue = "foobar"
)

func init() {
	os.Setenv("MY_TEST_VAR", envVarValue)
	config.addEnvVar("MY_TEST_VAR")
}

func TestEnvVarPopulation(t *testing.T) {
	ctx, _ := newContext(&config)

	if ctx.Env["MY_TEST_VAR"] != envVarValue {
		t.Errorf("test issue with '%s'", ctx.Env["MY_TEST_VAR"])
	}
}

func TestEnvVarPopulation_not_defined(t *testing.T) {
	myconfig := Config{}
	myconfig.addEnvVar("NOT_EXISTING_VAR")
	_, err := newContext(&myconfig)

	if err == nil {
		t.Error("Error expected for undefined env var")
	} else if !strings.Contains(err.Error(), "NOT_EXISTING_VAR") {
		t.Error("Error message does not contain NOT_EXISTING_VAR, value=", err.Error())
	}
}

func TestTemplate(t *testing.T) {
	testTplPath1 := "./tests/tmp/mytemplate1.txt"
	testTplPath2 := "./tests/tmp/mytemplate2.txt"
	testTplPath3 := "./tests/tmp/mytemplate3.txt"
	defer os.RemoveAll(filepath.Dir(testTplPath1))

	currentUser, err := user.Current()
	if err != nil {
		t.Errorf("could not lookup current user, error: %v", err)
	}

	currentGroup, err := lookupGroupById(currentUser.Gid)

	config.clearTemplates()
	config.addTemplate("./tests/test.tpl", testTplPath1, currentUser.Username, currentGroup.Name, "")
	config.addTemplate("./tests/test.tpl", testTplPath2, currentUser.Username, currentGroup.Name, "0600")
	config.addTemplate("./tests/test.tpl", testTplPath3, currentUser.Uid, currentGroup.Gid, "0750")
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
	assertFileExistsAndContains(testTplPath3, envVarValue, t)

	assertFileMode(testTplPath2, 0600, t)
	assertFileMode(testTplPath3, 0750, t)

	assertFileOwnerGroup(testTplPath1, currentUser.Uid, currentGroup.Gid, t)
	assertFileOwnerGroup(testTplPath2, currentUser.Uid, currentGroup.Gid, t)
	assertFileOwnerGroup(testTplPath3, currentUser.Uid, currentGroup.Gid, t)
}

func TestTemplate_notExistingOwnerGroup(t *testing.T) {
	testTplPath1 := "./tests/tmp/mytemplate1.txt"
	defer os.RemoveAll(filepath.Dir(testTplPath1))

	config.clearTemplates()
	config.addTemplate("./tests/test.tpl", testTplPath1, "notExistingUser", "notExistingGroup", "")
	ctx, _ := newContext(&config)

	err := generateTemplates(config.Templates, ctx)
	if err == nil {
		t.Error("error expected")
	} else if !strings.Contains(err.Error(), "notExistingUser") {
		t.Errorf("error should contain notExistingUser, error: %s", err.Error())
	}

	config.clearTemplates()
	config.addTemplate("./tests/test.tpl", testTplPath1, "", "notExistingGroup", "")
	ctx, _ = newContext(&config)

	err = generateTemplates(config.Templates, ctx)
	if err == nil {
		t.Error("error expected")
	} else if !strings.Contains(err.Error(), "notExistingGroup") {
		t.Errorf("error should contain notExistingGroup, error: %s", err.Error())
	}
}

func TestTemplate_invalidFileMode(t *testing.T) {
	testTplPath1 := "./tests/tmp/mytemplate1.txt"
	defer os.RemoveAll(filepath.Dir(testTplPath1))

	config.clearTemplates()
	config.addTemplate("./tests/test.tpl", testTplPath1, "", "", "99999")
	ctx, _ := newContext(&config)

	err := generateTemplates(config.Templates, ctx)
	if err == nil {
		t.Error("error expected")
	} else if !strings.Contains(err.Error(), "99999") {
		t.Errorf("error should contain 99999, error: %s", err.Error())
	}
}

func TestTemplate_chmodError(t *testing.T) {
	testTplPath1 := "./tests/tmp/mytemplate1.txt"
	defer os.RemoveAll(filepath.Dir(testTplPath1))

	config.clearTemplates()
	config.addTemplate("./tests/test.tpl", testTplPath1, "898989", "", "")
	ctx, _ := newContext(&config)

	err := generateTemplates(config.Templates, ctx)
	if err == nil {
		t.Error("error expected")
	} else if !strings.Contains(err.Error(), "898989") {
		t.Errorf("error should contain 898989, error: %s", err.Error())
	}
}

func TestMyMain(t *testing.T) {
	testTplPath1 := "./tests/dest/mytemplate.txt"
	defer os.RemoveAll(filepath.Dir(testTplPath1))

	os.Args = append(os.Args, "-config=tests/config.json")

	main()

	assertFileExistsAndContains(testTplPath1, envVarValue, t)
	assertFileMode(testTplPath1, 0700, t)
}

func TestIsValidFileMode(t *testing.T) {

	type fileModeTests struct {
		input string
		valid bool
	}

	var fileModes = []fileModeTests{
		fileModeTests{input: "0777", valid: true},
		fileModeTests{input: "0123", valid: true},
		fileModeTests{input: "0755", valid: true},
		fileModeTests{input: "2755", valid: true},
		fileModeTests{input: "7755", valid: true},
		fileModeTests{input: "8755", valid: false},
		fileModeTests{input: "", valid: false},
		fileModeTests{input: "1", valid: false},
		fileModeTests{input: "12", valid: false},
		fileModeTests{input: "123", valid: false},
		fileModeTests{input: "12345", valid: false},
	}

	for _, value := range fileModes {
		if isValidFileMode(value.input) != value.valid {
			t.Errorf("File mode %s should be %t", value.input, value.valid)
		}
	}
}

func assertFileExistsAndContains(file, content string, t *testing.T) {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		t.Errorf("file %s does not exist", file)
	} else {
		b, err := ioutil.ReadFile(file)
		if err != nil {
			t.Errorf("file %s could not be read, error: %v", file, err)
		}
		fileContent := string(b)
		if !strings.Contains(fileContent, content) {
			t.Errorf("File %s does not contain %s", file, content)
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

func assertFileOwnerGroup(file, expectedUserId, expectedGroupId string, t *testing.T) {
	fileStat, err := os.Stat(file)
	if err != nil {
		t.Errorf("File %s stat error %v", file, err)
	}

	uid := fileStat.Sys().(*syscall.Stat_t).Uid
	gid := fileStat.Sys().(*syscall.Stat_t).Gid

	userId := strconv.FormatUint(uint64(uid), 10)
	groupId := strconv.FormatUint(uint64(gid), 10)

	if expectedUserId != userId || expectedGroupId != groupId {
		t.Errorf("uid/gid of file %s is %s/%s but expected is %s/%s", file, userId, groupId, expectedUserId,
			expectedUserId)
	}
}
