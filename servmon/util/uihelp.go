package util

import (
	"log"

	"../data"
	ui "github.com/gizak/termui"
)

// UtilLogger is the logger for util
var UtilLogger *log.Logger

func linuxGaugeModifer(g *ui.Gauge, source data.LinuxDataSource) {
	g.Label = "test"
	value, err := source.GetMostRecentLoad(data.Min5)
	if err != nil {
		UtilLogger.Println("Couldn't get load because: ", err.Error())
	}
	g.Percent = int((value / 8) * 100)
}

func gaugeModifier(g *ui.Gauge, source data.DataSource) {
	for datapoint := range source.DataChan() {
		g.Label = datapoint.Metric
		// TODO: change this to number of processors
		g.Percent = int((datapoint.Value / float32(8)) * 100)
	}
}
