package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
)

var initialNumOfGoRoutines int

func isScenarioRunning() bool {
	numOfGoRoutines := runtime.NumGoroutine()
	log.Println("running GoRoutine:", numOfGoRoutines)
	return numOfGoRoutines > initialNumOfGoRoutines
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
			statusCodeCounts := make(map[int]int)
			var mu sync.Mutex

			for i := 0; i < requestCount; i++ {
				wg.Add(1)
				go func(id int) {
					defer wg.Done()
					code := post(url, id)
					mu.Lock()
					statusCodeCounts[code]++
					mu.Unlock()
				}(i)
			}
			wg.Wait()

			// log summary only after all get logs are printed
			wg.Add(1)
			go func() {
				defer wg.Done()
				fmt.Println("------Summary-------")
				printSummary(statusCodeCounts)
			}()
			wg.Wait()

		}
	}
}

func post(url string, id int) int {
	body := []byte(fmt.Sprintf("{\"id\":\"%d\"}", id))
	idStr := strconv.Itoa(id)
	start := time.Now()

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Fatalf("Error sending POST request to %s: %s", url, err)
	}
	defer resp.Body.Close()

	end := time.Now()
	elapsed := time.Since(start)
	log.Println("ID: "+ idStr + ", STATUS: " + resp.Status + ", START: " + start.Format("15:04:05") + ", END: " + end.Format("15:04:05") + ", TIME: " + elapsed.String())

	return resp.StatusCode
}

func printSummary(statusCodeCounts map[int]int) {
	for code, count := range statusCodeCounts {
		fmt.Printf("Requests that returned %d: %d\n", code, count)
	}
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
	initialNumOfGoRoutines = runtime.NumGoroutine() + 2
	fmt.Println("Initial number of go routines:", initialNumOfGoRoutines)
	err := http.ListenAndServe(":"+config.port, nil)
	if err != nil {
		log.Fatal(err)
	}

}
