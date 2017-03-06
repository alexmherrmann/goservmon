package data

import (
	"../util"
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"strings"
)

/*
A linux data source
Does not store passwords in memory long-term for safety reasons, creates the connection and holds onto it
*/
type LinuxDataSource struct {
	hostname   string
	client     *ssh.Client
	dataChan   chan DataPoint
	processors float32
}

// Close the channel and ssh client
func (l *LinuxDataSource) Close() {
	close(l.dataChan)
	l.client.Close()
}

const (
	// one minute average
	Min1 int = iota
	// / 5 minute average
	Min5
	// 10 minute average
	Min10
)

func (source *LinuxDataSource) GetNewSession() (*ssh.Session, error) {
	// TODO: error handling
	sesh, _ := source.client.NewSession()

	return sesh, nil

}

/*
TODO: implement which loadavg to select
*/
func (dataSource *LinuxDataSource) GetMostRecentLoad(avg int) float32 {

	//TODO: implement error handling
	newSesh, err := dataSource.client.NewSession()
	util.HandleError(err)

	readBytes := new(bytes.Buffer)
	newSesh.Stdout = readBytes

	//TODO: implement error handling
	err = newSesh.Run("/usr/bin/env cat /proc/loadavg")
	util.HandleError(err)

	log.Printf("Got [%s] for loadavg", strings.Trim(readBytes.String(), "\n"))

	var avg1, avg2, avg3 float32

	fmt.Sscanf(readBytes.String(), "%f %f %f", &avg1, &avg2, &avg3)

	log.Printf("got %f, %f, %f", avg1, avg2, avg3)

	switch avg {
	case Min1:
		return avg1

	default:
		fallthrough
	case Min5:
		return avg2

	case Min10:
		return avg3
	}
}

func (dataSource *LinuxDataSource) DataChan() chan DataPoint {
	return dataSource.dataChan
}

func (dataSource *LinuxDataSource) GetAllAvailablePoints() []DataPoint {
	panic("implement me")
}

// TODO: implement error handling for real
func GetNewLinuxDataSource(hostname string, username string, password string, processors int) (LinuxDataSource, error) {
	toReturn := LinuxDataSource{}

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{ssh.Password(password)},
	}

	client, err := ssh.Dial("tcp", hostname, config)
	util.HandleError(err)
	toReturn.client = client
	toReturn.processors = float32(processors)

	return toReturn, nil
}

func GetNewLinuxDataSource8Processors(hostname string, username string, password string) (LinuxDataSource, error) {
	return GetNewLinuxDataSource(hostname, username, password, 8)
}
