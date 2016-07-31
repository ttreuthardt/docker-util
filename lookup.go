package main

import (
	"fmt"
	"strconv"
	"unsafe"
)

/*
#include <grp.h>
#include <stdlib.h>
*/
import "C"

type Group struct {
	Gid  string
	Name string
}

func lookupGroupById(groupId string) (*Group, error) {
	gid, err := strconv.Atoi(groupId)
	if err != nil {
		return nil, fmt.Errorf("Given gid string could not be converted into an int: %v", err)
	}
	c_gid := C.gid_t(gid)
	c_grp, err := C.getgrgid(c_gid)
	if err != nil {
		return nil, fmt.Errorf("error while lookup group %s: %v", groupId, err)
	}

	if c_grp == nil {
		return nil, fmt.Errorf("Unknown group id %s", groupId)
	}

	return convertGroup(c_grp), nil
}

func lookupGroupByName(groupName string) (*Group, error) {
	c_groupName := C.CString(groupName)
	defer C.free(unsafe.Pointer(c_groupName))

	c_grp, err := C.getgrnam(c_groupName)
	if err != nil {
		return nil, fmt.Errorf("error while lookup group %s: %v", groupName, err)
	}

	if c_grp == nil {
		return nil, fmt.Errorf("Unknown group %s", groupName)
	}
	return convertGroup(c_grp), nil
}

func convertGroup(c_grp *C.struct_group) *Group {
	group := &Group{
		Gid:  strconv.Itoa(int(c_grp.gr_gid)),
		Name: C.GoString(c_grp.gr_name),
	}
	return group
}
