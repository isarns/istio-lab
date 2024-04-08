package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"slices"
	"strings"
	"sync"
	"time"
)

func isScenarioRunning() bool {
	numOfGoRoutine := runtime.NumGoroutine()
	log.Println("running GoRoutine:", numOfGoRoutine)
	return numOfGoRoutine > 3
}

func scenario(stopChannel chan bool, config config, path string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		log.Println(req.Method, req.RequestURI)
		if isScenarioRunning() {
			fmt.Fprintf(w, "Another scenario is running please use the /stop endpoint before running another")
			return
		}
		go ProcessData(
			config.serviceBUrl+path,
			config.timeToRun,
			config.requestCount,
			stopChannel,
		)
		fmt.Fprintf(w, "running: "+strings.Split(path, "/")[1])
	}
}

func stopScenario(stopChannel chan bool) http.HandlerFunc {
	//TODO: add check to see if scenario is running
	return func(w http.ResponseWriter, r *http.Request) {
		if !isScenarioRunning() {
			fmt.Fprintf(w, "No scenario is running please use the /scenario[A-C] endpoint before running /stop")
			return
		}
		log.Println("stopping...")
		fmt.Fprintf(w, "stopping...")
		stopChannel <- false
	}
}

func test(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "OK")
}

func logDetails(config config, appName string) {
	log.Println(appName, "runs on http://127.0.0.1:"+config.port, "\n",
		"config:", config)
}

func ProcessData(url string, timeToRun []int, requestCount int, stopChannel chan bool) {
	for {
		select {
		case <-stopChannel:
			log.Println("Stop go routine")
			return
		default:
			currentTime := time.Now()

			if !slices.Contains(timeToRun, currentTime.Second()) {
				time.Sleep(1 * time.Second)
				continue
			}

			fmt.Println("------Start-------")
			wg := sync.WaitGroup{}
			for i := 0; i < requestCount; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					get(url)
				}()
			}
			wg.Wait()
		}
	}
}

func get(url string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error: %s", err)
		return
	}
	defer resp.Body.Close()
	end := time.Now()
	elapsed := time.Since(start)

	log.Println("STATUS: " + resp.Status + ", START: " + getTime(start) + ", END: " + getTime(end) + ", TIME: " + elapsed.String())
}

func getTime(t time.Time) string {
	return strings.Split(t.String(), " ")[1]
}

func main() {
	config := initConfig()
	stopChannel := make(chan bool)
	logDetails(config, "App A")
	http.HandleFunc("/scenarioA", scenario(stopChannel, config, "/scenarioA"))
	http.HandleFunc("/scenarioB", scenario(stopChannel, config, "/scenarioB"))
	http.HandleFunc("/scenarioC", scenario(stopChannel, config, "/scenarioC"))
	http.HandleFunc("/stop", stopScenario(stopChannel))
	http.HandleFunc("/test", test)
	err := http.ListenAndServe(":"+config.port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
