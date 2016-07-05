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
		t.Fatalf("test issue with '%s'\n", ctx.Env["HOME"])
	}
}

func TestTemplate(t *testing.T) {
	config := Config{}
	config.addEnvVar("HOME")
	config.addTemplate("./tests/test.tpl")
	ctx := NewContext()

	populateContext(ctx, config)

	handleTemplate("./tests/test.tpl", "", ctx)
}

