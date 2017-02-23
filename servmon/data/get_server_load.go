package data

import (
	"golang.org/x/crypto/ssh"
	"../util"
	"bytes"
	"log"
	"strings"
	"fmt"
)

/*
A linux data source
Does not store passwords in memory long-term for safety reasons, creates the connection and holds onto it
 */
type LinuxDataSource struct {
	hostname string
	client *ssh.Client
}

type LinuxDataSourceConnectionError struct {
	originalMsg string
}

func (l LinuxDataSourceConnectionError) Error() string {
	return l.originalMsg
}


func (dataSource *LinuxDataSource) GetMostRecentLoad() float32 {

	//TODO: implement error handling
	newSesh, err := dataSource.client.NewSession()
	util.HandleError(err)

	readBytes := new(bytes.Buffer)
	newSesh.Stdout = readBytes

	//TODO: implement error handling
	err = newSesh.Run("/usr/bin/env cat /proc/loadavg")
	util.HandleError(err)

	log.Printf("Got [%s] for loadavg", strings.Trim(readBytes.String(), "\n"))

	var min1, min2, min3 float32

	fmt.Sscanf(readBytes.String(), "%f %f %f", &min1, &min2, &min3)

	log.Printf("got %f, %f, %f", min1, min2, min3)

	return 0.8
}

func (dataSource *LinuxDataSource) DataChan() (chan DataPoint) {
	panic("implement me")
}

func (dataSource *LinuxDataSource) GetAllAvailablePoints() []DataPoint {
	panic("implement me")
}

// TODO: implement error handling for real
func GetNewLinuxDataSource(hostname string, username string, password string) (LinuxDataSource, error) {
	toReturn := LinuxDataSource{}

	config := &ssh.ClientConfig{
		User:username,
		Auth: []ssh.AuthMethod { ssh.Password(password) },
	}

	client, err := ssh.Dial("tcp", hostname, config)
	util.HandleError(err)
	toReturn.client = client

	return toReturn, nil
}