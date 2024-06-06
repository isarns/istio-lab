package main

import (
	"github.com/isarns/IstioCircuitBreaker/utils"
	"log"
	"strconv"
)

type config struct {
	port        string
	timeToSleep int
	serviceBUrl string
	serviceCUrl string
}

func initConfig() config {
	port := utils.GetEnv("PORT", "9080")
	serviceBUrl := utils.GetEnv("SERVICE_B_URL", "http://127.0.0.1:9080")
	serviceCUrl := utils.GetEnv("SERVICE_C_URL", "http://127.0.0.1:9070")
	timeToSleep, err := strconv.Atoi(utils.GetEnv("TIME_TO_SLEEP", "1"))
	if err != nil {
		timeToSleep = 10
		log.Println("could not convert TIME_TO_SLEEP to int will use 1 as default")
	}
	config := config{
		port:        port,
		timeToSleep: timeToSleep,
		serviceBUrl: serviceBUrl,
		serviceCUrl: serviceCUrl,
	}
	return config
}
