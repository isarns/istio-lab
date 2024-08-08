package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"runtime"
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

func scenario(config config, baseUrl string, path string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		log.Println(req.Method, req.RequestURI)
		requestCount, delay := handleParams(req)
		go ProcessData(
			baseUrl+path,
			requestCount,
			delay,
		)
		fmt.Fprintf(w, "running scenario %s, url %s,  with requestCount: %d, delay: %s", strings.Split(path, "/")[1], baseUrl+path, requestCount, delay.String())
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

func ProcessData(url string, requestCount int, delay time.Duration) {
	fmt.Println("------Start-------")
	wg := sync.WaitGroup{}
	statusCodeCounts := make(map[int]int)
	var mu sync.Mutex

	for i := 0; i < requestCount; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			code := post(url, id, delay)
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

func post(url string, id int, delay time.Duration) int {
	body := []byte(fmt.Sprintf("{\"id\":\"%d\"}", id))
	idStr := strconv.Itoa(id)
	time.Sleep(delay * time.Duration(id))
	start := time.Now()

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Fatalf("Error sending POST request to %s: %s", url, err)
	}
	defer resp.Body.Close()

	end := time.Now()
	elapsed := time.Since(start)
	log.Println("ID: " + idStr + ", STATUS: " + resp.Status + ", START: " + start.Format("15:04:05") + ", END: " + end.Format("15:04:05") + ", TIME: " + elapsed.String())

	return resp.StatusCode
}

func printSummary(statusCodeCounts map[int]int) {
	for code, count := range statusCodeCounts {
		fmt.Printf("Requests that returned %d: %d\n", code, count)
	}
}

func handleRequestCountParam(req *http.Request) int {
	var requestCount int
	requestCountFromParams := req.URL.Query()["requestCount"]
	if len(requestCountFromParams) > 0 {
		tempRequestCount, err := strconv.Atoi(requestCountFromParams[0])
		if err != nil {
			log.Println("could not convert requestCount to int will use 20 as default")
		}
		requestCount = tempRequestCount
	} else {
		log.Println("requestCount not provided, using default value of 20")
		requestCount = 20
	}
	return requestCount

}

func handleDelayParam(req *http.Request) time.Duration {
	var delay time.Duration
	delayFromParams := req.URL.Query()["delay"]
	if len(delayFromParams) > 0 {
		delayDuratin, err := time.ParseDuration(delayFromParams[0])
		if err != nil {
			log.Println("could not convert delay to time.Duration will use 0s as default")
			delay = 0 * time.Second
		} else {
			delay = delayDuratin
		}
	} else {
		delay = 0 * time.Second
	}
	return delay
}

func handleParams(req *http.Request) (int, time.Duration) {
	requestCount := handleRequestCountParam(req)
	delay := handleDelayParam(req)
	return requestCount, delay
}

func withParams(w http.ResponseWriter, req *http.Request) {
	requestCount, Delay := handleParams(req)
	fmt.Fprintf(w, "requestCount: %d, delay: %s", requestCount, Delay.String())
}

func main() {
	config := initConfig()
	stopChannel := make(chan bool)
	logDetails(config, "App A")
	http.HandleFunc("/withParams", withParams)
	http.HandleFunc("/scenarioA", scenario(config, config.serviceBUrl, "/scenarioA"))
	http.HandleFunc("/scenarioB", scenario(config, config.serviceBUrl, "/scenarioB"))
	http.HandleFunc("/scenarioC", scenario(config, config.serviceBUrl, "/scenarioC"))
	http.HandleFunc("/scenarioD", scenario(config, config.serviceCUrl, "/scenarioD"))
	http.HandleFunc("/stop", stopScenario(stopChannel))
	http.HandleFunc("/test", test)
	err := http.ListenAndServe(":"+config.port, nil)
	if err != nil {
		log.Fatal(err)
	}

}
