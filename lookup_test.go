package main

import (
	"os/user"
	"strings"
	"testing"
)

var currentUser, _ = user.Current()

func TestLookupGroup(t *testing.T) {
	groupById, err := LookupGroupById(currentUser.Gid)
	if err != nil {
		t.Error(err.Error())
	}

	groupByName, err := LookupGroupByName(groupById.Name)
	if err != nil {
		t.Error(err.Error())
	}

	if groupById.Name != groupByName.Name {
		t.Errorf("the lookups did not return the same group, %v != %v", groupById, groupByName)
	}
}

func TestLookupGroupById_errors(t *testing.T) {
	_, err := LookupGroupById("ddddd")
	if err == nil {
		t.Error("error expected")
	}

	_, err = LookupGroupById("1010101010")
	if err == nil {
		t.Error("error expected")
	}
}

func TestLookupGroupByName_errors(t *testing.T) {
	groupName := "notExistingGroupName"
	_, err := LookupGroupByName(groupName)
	if err == nil {
		t.Error("error expected")
	}
	if !strings.Contains(err.Error(), groupName) {
		t.Errorf("error should contain group name, error: %s", err.Error())
	}

	_, err = LookupGroupByName("")
	if err == nil {
		t.Error("error expected")
	}
}
