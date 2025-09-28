package config

import (
	"log"
	"os"
	"runtime"
	"strconv"
)

func SetupResources() {

	//LIMITING CPU
	cpuEnv := os.Getenv("MAX_VCPU")
	numCPU := runtime.NumCPU()
	if cpuEnv != "" {
		if cpuVal, err := strconv.Atoi(cpuEnv); err == nil && cpuVal > 0 {
			numCPU = cpuVal
		}
	}
	log.Printf("Setting CPU cores to %d", numCPU)
	runtime.GOMAXPROCS(numCPU)

}
