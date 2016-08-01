package main

type Group struct {
	Gid  string
	Name string
}

func LookupGroupById(groupId string) (*Group, error) {
	return lookupGroupById(groupId)
}

func LookupGroupByName(groupName string) (*Group, error) {
	return lookupGroupByName(groupName)
}
