package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func addSleep(next http.HandlerFunc, sleepTime time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		time.Sleep(sleepTime)
		next.ServeHTTP(w, req)
	}
}

func addLog(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		log.Println(req.Method, req.URL.Path)
		next.ServeHTTP(w, req)
	}
}

func scenarioA(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "scenarioA\n")
}

func scenarioB(config config) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		makeGetRequest("http://127.0.0.1:" + config.port + "/talkingToMyself")
	}
}

func scenarioC(config config) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		makeGetRequest(config.serviceCUrl + "/scenarioC")
	}
}

func makeGetRequest(url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error: %s", err)
		return
	}
	defer resp.Body.Close()
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
	http.HandleFunc("/scenarioA", addLog(addSleep(scenarioA, timeToSleep)))
	http.HandleFunc("/scenarioB", addLog(scenarioB(config)))
	http.HandleFunc("/scenarioC", addLog(scenarioC(config)))
	http.HandleFunc("/talkingToMyself", addLog(addSleep(talkingToMyself, timeToSleep)))
	http.HandleFunc("/test", test)
	err := http.ListenAndServe(":"+config.port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
