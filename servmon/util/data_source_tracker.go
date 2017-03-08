package util

import (
	"../data"
	ui "github.com/gizak/termui"
)

// TrackerPair is a simple class matching a datasource to a gauge... for now
type TrackerPair struct {
	Source *data.LinuxDataSource
	Gauge  *ui.Gauge
}

var tracker struct {
	Pairs []TrackerPair
	rows []*ui.Row
}

// RegisterServer takes a LinuxDataSource and adds it to the track list
func RegisterServer(source *data.LinuxDataSource) {
	pair := TrackerPair{Source: source, Gauge: ui.NewGauge()}
	pairs := append(tracker.Pairs, pair)
	tracker.Pairs = pairs
	tracker.rows = append(tracker.rows, ui.NewRow(
		ui.NewCol(12, 0, pair.Gauge),
	))
}

// GetPairs just returns all of the TrackerPairs
func GetPairs() []TrackerPair {
	return tracker.Pairs
}

func GetRows() []*ui.Row {
	return tracker.rows
}