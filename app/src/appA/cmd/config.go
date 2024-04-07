package main

import (
	"log"
	"os"
	"strconv"
	"strings"
)

type config struct {
	port         string
	timeToRun    []int
	requestCount int
	serviceBUrl  string
}

func initConfig() config {
	port := getEnv("PORT", "9090")
	timeToRun := fromStringToIntArray(getEnv("TIME_TO_RUN", "00,30"))
	serviceBUrl := getEnv("SERVICE_B_URL", "http://127.0.0.1:9080")
	requestCount, err := strconv.Atoi(getEnv("REQUEST_COUNT", "20"))
	if err != nil {
		requestCount = 20
		log.Println("could not convert REQUEST_COUNT to int will use 20 as default")
	}
	config := config{
		port:         port,
		timeToRun:    timeToRun,
		requestCount: requestCount,
		serviceBUrl:  serviceBUrl,
	}
	return config
}

func fromStringToIntArray(s string) []int {
	var result []int
	for _, v := range strings.Split(s, ",") {
		s, err := strconv.Atoi(v)
		if err != nil {
			log.Panic("cant convert", v, "to string.")
		}
		result = append(result, s)
	}
	return result
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
