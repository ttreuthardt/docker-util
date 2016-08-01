package main

import (
	"os/user"
	"testing"
)

func TestLookupGroup(t *testing.T) {
	currentUser, _ := user.Current()
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
