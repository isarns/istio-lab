package main

import (
	"log"
	"os"
	"strconv"
)

type config struct {
	port        string
	timeToSleep int
}

func initConfig() config {
	port := getEnv("PORT", "9070")
	timeToSleep, err := strconv.Atoi(getEnv("TIME_TO_SLEEP", "10"))
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

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
