package main

import (
	"strings"
	"testing"
)

func TestReadConfigFile(t *testing.T) {
	config, err := readConfig("tests/test_config.json")
	if err != nil {
		t.Errorf("readConfig, %v", err)
	}

	if len(config.Envvars) != 3 {
		t.Error("Env variables not available")
	}

	if len(config.Templates) != 2 {
		t.Errorf("2 Templates expected got only %v", len(config.Templates))
	}

	for _, tpl := range config.Templates {

		if strings.HasSuffix(tpl.DestPath, "tpl1") {
			if tpl.FileMode != "0777" {
				t.Errorf("file mode for tpl1 is not as expected %v", tpl.FileMode)
			}
		}

		if strings.HasSuffix(tpl.DestPath, "tpl2") {
			if tpl.FileMode != "" {
				t.Errorf("file mode for tpl2 is not as expected %v", tpl.FileMode)
			}
		}

	}

}
