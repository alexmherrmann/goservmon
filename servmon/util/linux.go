package util

import (
	"bytes"
	"errors"
	"strings"

	"golang.org/x/crypto/ssh"
)

// ErrGetNumberOfCores means that for whatever reason, we couldn't get the core count
var ErrGetNumberOfCores = errors.New("Couldn't get the core count for ssh cores")

// GetNumberOfCores takes an ssh session and tries to intuit the number of cores on the system
func GetNumberOfCores(session ssh.Session) (int, error) {
	readBytes := new(bytes.Buffer)
	session.Stdout = readBytes

	var err error
	// TODO: error handling
	err = session.Run("/usr/bin/env cat /proc/cpuinfo")
	if err != nil {
		return -1, ErrGetNumberOfCores
	}

	count := GetNumberOfCoresFromCPUInfo(readBytes)

	return count, nil
}

// GetNumberOfCoresFromCPUInfo takes some bytes, parses it and returns the number of times "processor" appears in it
func GetNumberOfCoresFromCPUInfo(readBytes *bytes.Buffer) int {
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
