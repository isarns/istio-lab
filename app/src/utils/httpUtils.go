package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type IDStruct struct {
	ID string `json:"id"`
}

func countForSeconds(sleepTime time.Duration) {
	endTime := time.Now().Add(sleepTime)
	count := 0
	for time.Now().Before(endTime) {
		count++
		time.Sleep(10 * time.Microsecond)
	}
}

func WithAddSleep(next http.HandlerFunc, sleepTime time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		countForSeconds(sleepTime)
		next.ServeHTTP(w, req)
	}
}

func WithAddLog(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		bodyBytes := ReadBody(req)
		log.Println(req.Method, req.URL.Path, formatID(bodyBytes))
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		next.ServeHTTP(w, req)
	}
}

func MakeGetRequest(url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error: %s", err)
		return
	}
	defer resp.Body.Close()
}

func MakePostRequest(url string, body []byte ) {
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Fatalf("Error sending POST request to %s: %s", url, err)
	}
	defer resp.Body.Close()
}

func ReadBody(req *http.Request) []byte {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	return body
}

// formatID takes a JSON string and returns a string in the format "id=value"
func formatID(bodyBytes []byte) string {
	var idStruct IDStruct
	err := json.Unmarshal(bodyBytes, &idStruct)
	if err != nil {
		return "id=NONE"
	}
	return fmt.Sprintf("id=%s", idStruct.ID)
}