package util

import (
	"bytes"
	"golang.org/x/crypto/ssh"
	"strings"
)

func GetNumberOfCores(session ssh.Session) int {
	readBytes := bytes.Buffer{}
	session.Stdout = readBytes

	// TODO: error handling
	_ = session.Run("/usr/bin/env cat /proc/cpuinfo")

	count := 0
	for _, str := range strings.Split(readBytes.String(), "\n") {
		if str[:8] == "processor" {
			count++
		}
	}

	return count
}
