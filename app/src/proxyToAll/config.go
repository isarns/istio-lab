package main

import (
	"github.com/isarns/IstioCircuitBreaker/utils"
	"log"
	"strconv"
)

type config struct {
	servicePort          int
	serviceLabelSelector string
	serviceNamespace     string
}

func initConfig() config {
	serviceLabelSelector := utils.GetEnv("SERVICE_LABEL_SELECTOR", "app=app-a")
	serviceNamespace := utils.GetEnv("SERVICE_NAMESPACE", "apps")
	servicePort, err := strconv.Atoi(utils.GetEnv("SERVICE_PORT", "8080"))
	if err != nil {
		servicePort = 8080
		log.Println("could not convert SERVICE_PORT to int will use 8080 as default")
	}
	config := config{
		servicePort:          servicePort,
		serviceLabelSelector: serviceLabelSelector,
		serviceNamespace:     serviceNamespace,
	}
	return config
}
