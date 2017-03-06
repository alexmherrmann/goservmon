package util

import (
	"bytes"
	"golang.org/x/crypto/ssh"
	"strings"
)


//func GetNumberOfCoresFromCpuinfo(input string) int {
//
//}

func GetNumberOfCores(session ssh.Session) int {
	readBytes := new(bytes.Buffer)
	session.Stdout = readBytes

	// TODO: error handling
	_ = session.Run("/usr/bin/env cat /proc/cpuinfo")

	count := GetNumberOfCoresFromCpuInfo(readBytes)

	return count
}
func GetNumberOfCoresFromCpuInfo(readBytes *bytes.Buffer) int {
	count := 0
	for _, str := range strings.Split(readBytes.String(), "\n") {
		if len(str) >= 8 {
			if str[:9] == "processor" {
				count++
			}
		} else {
			continue
		}
	}
	return count
}

