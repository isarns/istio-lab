package main

import (
	"log"
	"os"
	"strconv"
)

type config struct {
	port        string
	timeToSleep int
	serviceBUrl string
	serviceCUrl string
}

func initConfig() config {
	port := getEnv("PORT", "9080")
	serviceBUrl := getEnv("SERVICE_B_URL", "http://127.0.0.1:9080")
	serviceCUrl := getEnv("SERVICE_C_URL", "http://127.0.0.1:9070")
	timeToSleep, err := strconv.Atoi(getEnv("TIME_TO_SLEEP", "10"))
	if err != nil {
		timeToSleep = 10
		log.Println("could not convert TIME_TO_SLEEP to int will use 10 as default")
	}
	config := config{
		port:        port,
		timeToSleep: timeToSleep,
		serviceBUrl: serviceBUrl,
		serviceCUrl: serviceCUrl,
	}
	return config
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
