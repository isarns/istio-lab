package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/isarns/IstioCircuitBreaker/utils"
)

func scenarioC(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "scenarioC\n")
}

func test(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "OK")
}

func logDetails(config config, appName string) {
	log.Println(appName, "runs on http://127.0.0.1:"+config.port, "\n",
		"config:", config)
}

func main() {
	log.Println("Starting our simple http server.")
	config := initConfig()
	timeToSleep := time.Duration(config.timeToSleep) * time.Second
	logDetails(config, "App C")
	http.HandleFunc("/scenarioC", utils.WithLog(utils.WithSleep(scenarioC, timeToSleep)))
	http.HandleFunc("/test", test)
	err := http.ListenAndServe(":"+config.port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
