// +build darwin dragonfly freebsd !android,linux netbsd openbsd solaris
// +build cgo

package main

import (
	"fmt"
	"strconv"
	"unsafe"
)

/*
#include <grp.h>
#include <sys/types.h>
#include <stdlib.h>

static struct group * wrapper_getgrgid(int gid) {
 	return getgrgid(gid);
}
*/
import "C"

func lookupGroupById(groupId string) (*Group, error) {
	gid, err := strconv.Atoi(groupId)
	if err != nil {
		return nil, fmt.Errorf("Given gid string could not be converted into an int: %v", err)
	}
	/* wrapper_getgrgid has to be used to avoid C.gid_t as it does not work in linux */
	c_grp, _ := C.wrapper_getgrgid(C.int(gid))

	if c_grp == nil {
		return nil, fmt.Errorf("Unknown group id %s", groupId)
	}

	return convertGroup(c_grp), nil
}

func lookupGroupByName(groupName string) (*Group, error) {
	c_groupName := C.CString(groupName)
	defer C.free(unsafe.Pointer(c_groupName))

	c_grp, _ := C.getgrnam(c_groupName)

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
