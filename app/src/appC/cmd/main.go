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
	http.HandleFunc("/scenarioC", addLog(addSleep(scenarioC, timeToSleep)))
	http.HandleFunc("/test", test)
	err := http.ListenAndServe(":"+config.port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
