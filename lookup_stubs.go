// +build !cgo,!windows,!plan9 android

package main

import (
	"fmt"
)

func lookupGroupById(groupId string) (*Group, error) {
	return nil, fmt.Errorf("lookupGroupById not implemented")
}

func lookupGroupByName(groupName string) (*Group, error) {
	return nil, fmt.Errorf("lookupGroupByName not implemented")
}
