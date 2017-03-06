package main

import (
	"./servmon"
	"./servmon/data"
	"./servmon/util"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"syscall"

	ui "github.com/gizak/termui"
)

//func newGauge(source data.DataSource) ui.Gauge {
//}

func linuxGaugeModifer(g *ui.Gauge, source data.LinuxDataSource) {
	g.Label = "test"
	g.Percent = int( (source.GetMostRecentLoad(data.Min5)/8) * 100)
}

func gaugeModifier(g *ui.Gauge, source data.DataSource) {
	for datapoint := range source.DataChan() {
		g.Label = datapoint.Metric
		// TODO: change this to number of processors
		g.Percent = int((datapoint.Value / float32(8))*100)
	}
}

func main() {
	file, _ := os.OpenFile("servmon.log", os.O_CREATE | os.O_TRUNC | os.O_RDWR, 0755)
	log.SetOutput(file)
	settingsFile, err := os.Open("settings.yaml")
	util.HandleError(err)

	settingsBytes, err := ioutil.ReadAll(settingsFile)
	util.HandleError(err)

	config := servmon.Configuration{}
	yaml.Unmarshal(settingsBytes, &config)


	if false {
		fmt.Printf("Please input password for %s:", config.Username)
		passwdBytes, err := terminal.ReadPassword(int(syscall.Stdin))
		config.Password = string(passwdBytes)
		fmt.Println()
		util.HandleError(err)
	}

	serverToGet := config.Servers[0]
	log.Printf("Getting CPU load for %s\n", serverToGet)

	serverDataSource, err := data.GetNewLinuxDataSource8Processors(serverToGet, config.Username, config.Password)
	util.HandleError(err)
	load := serverDataSource.GetMostRecentLoad(data.Min5)

	log.Printf("Got load %f\n", load)

	ui.Init()
	g := ui.NewGauge()
	g.Height = 3
	g.Width = 50
	g.BarColor = ui.ColorRed
	g.BorderFg = ui.ColorWhite
	g.BorderLabelFg = ui.ColorCyan

	//ui.Handle("/timer/1s", func(e ui.Event) {
	ui.Handle("/sys/kbd/r", func(e ui.Event) {
		linuxGaugeModifer(g, serverDataSource)
		log.Println("rendering...")
		ui.Render(g)
	})

	ui.Handle("/sys/kbd/q", func(ui.Event) {
		// press q to quit
		ui.StopLoop()
	})


	g.Label = "ayylmao"
	g.Percent = 59

	defer ui.Close()
	ui.Render(g)
	ui.Loop()


}
