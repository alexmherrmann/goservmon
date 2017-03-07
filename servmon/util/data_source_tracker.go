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
}

// RegisterServer takes a LinuxDataSource and adds it to the track list
func RegisterServer(source *data.LinuxDataSource) {
	pairs := append(tracker.Pairs, TrackerPair{Source: source, Gauge: ui.NewGauge()})
	tracker.Pairs = pairs
}

// GetPairs just returns all of the TrackerPairs
func GetPairs() []TrackerPair {
	return tracker.Pairs
}
