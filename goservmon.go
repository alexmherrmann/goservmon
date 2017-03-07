package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"syscall"

	"./servmon"
	"./servmon/data"
	"./servmon/util"
	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/yaml.v2"

	ui "github.com/gizak/termui"
)

func setUpLogging(input *log.Logger) {
	file, err := os.OpenFile("servmon.log", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0755)
	if err != nil {
		panic("problem opening log file \"servmon.log\"")
	}
	input.SetOutput(file)

}

func main() {
	log := new(log.Logger)
	setUpLogging(log)
	data.DataLogger = log
	util.UtilLogger = log

	settingsFile, err := os.Open("settings.yaml")
	if err != nil {
		panic("couldn't open settings.yaml")
	}

	settingsBytes, err := ioutil.ReadAll(settingsFile)
	if err != nil {
		panic("error reading contents from settings.yaml")
	}

	config := servmon.Configuration{}
	yaml.Unmarshal(settingsBytes, &config)

	if false {
		fmt.Printf("Please input password for %s:", config.Username)
		passwdBytes, err := terminal.ReadPassword(int(syscall.Stdin))
		config.Password = string(passwdBytes)
		fmt.Println()
		util.HandleError(err)
	}

	ui.Init()

	// TODO: Put a big old loop here to keep track of data sources and gauges

}
