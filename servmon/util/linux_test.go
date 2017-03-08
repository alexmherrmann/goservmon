package util

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestGetNumberOfCores(t *testing.T) {
	load, _ := os.Open("/proc/cpuinfo")
	readbytes, _ := ioutil.ReadAll(load)
	toparse := bytes.NewBuffer(readbytes)
	fmt.Printf("processors: %d\n", GetNumberOfCoresFromCpuInfo(toparse))
}
