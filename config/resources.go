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

	//LIMITING MEMORY
	memEnv := os.Getenv("MAX_MEMORY_MB")
	maxMemoryMB := 0
	if memEnv != "" {
		if memVal, err := strconv.Atoi(memEnv); err == nil && memVal > 0 {
			maxMemoryMB = memVal
		}
	}
	if maxMemoryMB > 0 {
		log.Printf("Setting memory limit for buffers to %d MB", maxMemoryMB)
	} else {
		log.Println("No memory limit specified in env, using dynamic allocation.")
	}

}
