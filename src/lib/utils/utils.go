package utils

import (
	"errors"
	"io"
	"strings"

	"golang.org/x/sys/windows/registry"
)

const (
	// LxssBaseRoot is LOCAL_MACHINE
	LxssBaseRoot = registry.CURRENT_USER
	// LxssBaseKey is path of lxss registry
	LxssBaseKey = "Software\\Microsoft\\Windows\\CurrentVersion\\Lxss"
)

// DQEscapeString is escape string with double quote
func DQEscapeString(str string) string {
	if strings.Contains(str, " ") {
		str = strings.Replace(str, "\"", "\\\"", -1)
		str = "\"" + str + "\""
	}
	return str
}

//WslGetUUID gets distro guid key
func WslGetUUID(distributionName string) (uuid string, err error) {
	uuidList, tmpErr := WslGetUUIDList()
	if tmpErr != nil {
		err = tmpErr
		return
	}

	errStr := ""
	for _, loopUUID := range uuidList {
		key, loopErr := registry.OpenKey(LxssBaseRoot, LxssBaseKey+"\\"+loopUUID, registry.READ)
		if loopErr == nil || loopErr == io.EOF {
			str, _, itemErr := key.GetStringValue("DistributionName")
			if itemErr == nil || itemErr == io.EOF {
				if strings.EqualFold(str, distributionName) {
					uuid = loopUUID
					return
				}
			} else {
				errStr += "\n" + "    " + loopUUID + ":" + itemErr.Error()
			}
		} else {
			errStr += "\n" + "    " + loopUUID + ":" + loopErr.Error()
		}
	}
	err = errors.New("Registry Key Not found\n" + errStr)

	return
}

//WslGetUUIDList gets guid key lists
func WslGetUUIDList() (uuidList []string, err error) {
	baseKey, tmpErr := registry.OpenKey(LxssBaseRoot, LxssBaseKey, registry.READ)
	if tmpErr != nil && tmpErr != io.EOF {
		err = tmpErr
		return
	}
	uuidList, tmpErr = baseKey.ReadSubKeyNames(1024)
	if tmpErr != nil && tmpErr != io.EOF {
		err = tmpErr
		return
	}
	return
}
