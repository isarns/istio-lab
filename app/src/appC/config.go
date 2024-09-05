package main

import (
	"github.com/isarns/IstioCircuitBreaker/utils"
	"log"
	"strconv"
)

type config struct {
	port        string
	timeToSleep int
}

func initConfig() config {
	port := utils.GetEnv("PORT", "9070")
	timeToSleep, err := strconv.Atoi(utils.GetEnv("TIME_TO_SLEEP", "10"))
	if err != nil {
		timeToSleep = 10
		log.Println("could not convert TIME_TO_SLEEP to int will use 10 as default")
	}
	config := config{
		port:        port,
		timeToSleep: timeToSleep,
	}
	return config
}
