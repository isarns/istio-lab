package main

import (
	"github.com/isarns/IstioCircuitBreaker/utils"
	"log"
	"strconv"
)

type config struct {
	port         string
	timeToRun    []int
	requestCount int
	serviceBUrl  string
	serviceCUrl  string
}

func initConfig() config {
	port := utils.GetEnv("PORT", "9090")
	timeToRun := utils.FromStringToIntArray(utils.GetEnv("TIME_TO_RUN", "00,10,20,30,40,50"))
	serviceBUrl := utils.GetEnv("SERVICE_B_URL", "http://127.0.0.1:9080")
	serviceCUrl := utils.GetEnv("SERVICE_C_URL", "http://127.0.0.1:9070")
	requestCount, err := strconv.Atoi(utils.GetEnv("REQUEST_COUNT", "20"))
	if err != nil {
		requestCount = 20
		log.Println("could not convert REQUEST_COUNT to int will use 20 as default")
	}
	config := config{
		port:         port,
		timeToRun:    timeToRun,
		requestCount: requestCount,
		serviceBUrl:  serviceBUrl,
		serviceCUrl:  serviceCUrl,
	}
	return config
}
