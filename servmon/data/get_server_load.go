package data

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"strings"

	"golang.org/x/crypto/ssh"
)

/*LinuxDataSource is a DataSource for linux systems
Does not store passwords in memory long-term for safety reasons, creates the connection and holds onto it
*/
type LinuxDataSource struct {
	hostname   string
	client     *ssh.Client
	dataChan   chan DataPoint
	processors float32
}

// Close the channel and ssh client
func (dataSource *LinuxDataSource) Close() {
	close(dataSource.dataChan)
	dataSource.client.Close()
}

const (
	// Min1 is one minute average
	Min1 int = iota
	// Min5 is 5 minute average
	Min5
	// Min10 is 10 minute average
	Min10
)

// ErrGetNewSession is returned when GetNewSession could not complete the request for a session
var ErrGetNewSession = errors.New("Could not get session")

// GetNewSession allows you to get an ssh session from the LinuxDataSource
func (dataSource *LinuxDataSource) GetNewSession() (*ssh.Session, error) {
	// TODO: error handling
	sesh, err := dataSource.client.NewSession()
	if err != nil {
		return nil, ErrGetNewSession
	}
	return sesh, nil

}

// ErrGetLoad is returned when the Load could not get retrieved for an instance
var ErrGetLoad = errors.New("Could not get load")

/*
GetMostRecentLoad gets the most recent load avg for the system configured
*/
func (dataSource *LinuxDataSource) GetMostRecentLoad(avg int) (float32, error) {

	//TODO: implement error handling
	newSesh, err := dataSource.client.NewSession()
	if err != nil {
		return -1, ErrGetLoad
	}

	readBytes := new(bytes.Buffer)
	newSesh.Stdout = readBytes

	//TODO: implement error handling
	err = newSesh.Run("/usr/bin/env cat /proc/loadavg")
	if err != nil {
		return -1, ErrGetLoad
	}

	log.Printf("Got [%s] for loadavg", strings.Trim(readBytes.String(), "\n"))

	var avg1, avg2, avg3 float32

	fmt.Sscanf(readBytes.String(), "%f %f %f", &avg1, &avg2, &avg3)

	log.Printf("got %f, %f, %f", avg1, avg2, avg3)

	switch avg {
	case Min1:
		return avg1, nil

	default:
		fallthrough
	case Min5:
		return avg2, nil

	case Min10:
		return avg3, nil
	}
}

// DataChan gets the available channel, although it is not used yet
func (dataSource *LinuxDataSource) DataChan() chan DataPoint {
	panic("implement me")
	// return dataSource.dataChan
}

// GetAllAvailablePoints will return all of the points gathered so far
func (dataSource *LinuxDataSource) GetAllAvailablePoints() []DataPoint {
	panic("implement me")
}

// ErrCouldNotDialSSH is returned when an ssh connection could not be established
var ErrCouldNotDialSSH = errors.New("Could not dial the ssh server")

/*
GetNewLinuxDataSource creates a new LinuxDataSource
TODO: implement error handling for real
*/
func GetNewLinuxDataSource(hostname string, username string, password string, processors int) (*LinuxDataSource, error) {
	toReturn := new(LinuxDataSource)

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{ssh.Password(password)},
	}

	client, err := ssh.Dial("tcp", hostname, config)
	if err != nil {
		return nil, ErrCouldNotDialSSH
	}
	toReturn.client = client
	toReturn.processors = float32(processors)

	return toReturn, nil
}

/*
GetNewLinuxDataSource8Processors is a helper to just get a new linux data source configured
for use with 8 processors, it is defunct and should not be used
*/
func GetNewLinuxDataSource8Processors(hostname string, username string, password string) (*LinuxDataSource, error) {
	return GetNewLinuxDataSource(hostname, username, password, 8)
}
