package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/isarns/IstioCircuitBreaker/utils"
)

func scenarioA(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "scenarioA\n")
}

func scenarioB(config config) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		statusCode, body := utils.MakePostRequest(config.serviceCUrl+"/scenarioC", utils.ReadBody(req))
		w.WriteHeader(statusCode)
		w.Write(body)
	}
}

func scenarioC(config config) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		statusCode, body := utils.MakePostRequest(config.serviceCUrl+"/scenarioC", utils.ReadBody(req))
		w.WriteHeader(statusCode)
		w.Write(body)
	}
}


func talkingToMyself(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Much calculations\n")
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
	logDetails(config, "App B")
	http.HandleFunc("/scenarioA", utils.WithLog(utils.WithSleep(scenarioA, timeToSleep)))
	http.HandleFunc("/scenarioB", utils.WithLog(scenarioB(config)))
	http.HandleFunc("/scenarioC", utils.WithLog(scenarioC(config)))
	http.HandleFunc("/talkingToMyself", utils.WithLog(utils.WithSleep(talkingToMyself, timeToSleep)))
	http.HandleFunc("/test", test)
	err := http.ListenAndServe(":"+config.port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
