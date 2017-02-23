package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"./servmon"
	"./servmon/data"
	"golang.org/x/crypto/ssh/terminal"
	"syscall"
	"fmt"
	"log"
	"./servmon/util"
)

func main() {
	settingsFile, err := os.Open("settings.yaml")
	util.HandleError(err)

	settingsBytes, err := ioutil.ReadAll(settingsFile)
	util.HandleError(err)

	config := servmon.Configuration{}
	yaml.Unmarshal(settingsBytes, &config)

	fmt.Printf("Please input password for %s:", config.Username)
	passwdBytes, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	util.HandleError(err)
	config.Password = string(passwdBytes)

	serverToGet := config.Servers[0]
	log.Printf("Getting CPU load for %s\n", serverToGet)

	serverDataSource, err := data.GetNewLinuxDataSource(serverToGet, config.Username, config.Password)
	util.HandleError(err)
	load := serverDataSource.GetMostRecentLoad()

	log.Printf("Got load %f\n", load)


}