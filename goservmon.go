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
	input.SetFlags(log.Lshortfile | log.Ltime)

}

func main() {
	Logga := new(log.Logger)
	setUpLogging(Logga)
	data.DataLogger = Logga
	util.UtilLogger = Logga

	settingsFile, err := os.Open("settings.yaml")
	if err != nil {

		Logga.Fatalln("couldn't open settings.yaml")
	}

	settingsBytes, err := ioutil.ReadAll(settingsFile)
	if err != nil {
		Logga.Fatalln("error reading contents from settings.yaml")
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

	err = ui.Init()
	if err != nil {
		Logga.Fatalln(err.Error())
		os.Exit(-1)
	}
	defer ui.Close()

	// TODO: Put a big old loop here to keep track of data sources and gauges

	for _, server := range config.Servers {
		//TODO: Figure out how many processors!
		datasource, err := data.GetNewLinuxDataSource(server, config.Username, config.Password, 8)
		if err == nil {
			Logga.Println("Registering server ", server)
			util.RegisterServer(datasource)
		} else {
			Logga.Println("Error adding datasource for ", server)
		}
	}

	ui.Handle("/timer/5s", func(event ui.Event) {
		for _, pair := range util.GetPairs() {
			Logga.Println("modifying gauge")
			util.LinuxGaugeModifer(pair.Gauge, pair.Source)
		}
	})

	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})

	Logga.Println("beginning render")
	ui.Render(ui.NewGrid(util.GetRows()...))
	ui.Loop()

}
