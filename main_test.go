package main

import (
	"testing"
)

func TestEnvVarPopulation(t *testing.T) {
	config := Config{}
	config.addEnvVar("HOME")

	ctx := NewContext()

	populateContext(ctx, config)

	if len(ctx.Env["HOME"]) == 0 {
		t.Fatalf("test issue with '%s'", ctx.Env["HOME"])
	} else {
		t.Logf("env var HOME set to '%s'", ctx.Env["HOME"])
	}

}
