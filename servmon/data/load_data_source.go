package data


type DataPoint struct {
	value float32
}

type DataSource interface {
	DataChan() (chan DataPoint)
	GetAllAvailablePoints() []DataPoint
}