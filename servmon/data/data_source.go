package data

type DataPoint struct {
	Value  float32
	Metric string
}

type DataSource interface {
	DataChan() chan DataPoint
	GetAllAvailablePoints() []DataPoint
	Close()
}
