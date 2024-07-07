package utils

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"time"
)

type IDStruct struct {
	ID string `json:"id"`
}

// generateRandomString creates a random string of a given length.
// If no length is specified, it defaults to 20.
func generateRandomString(length ...int) string {
	l := 20 // Default length
	if len(length) > 0 {
		l = length[0]
	}

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, l)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			panic(err) // Handling error by panicking, can be replaced with better error handling
		}
		b[i] = charset[n.Int64()]
	}
	return string(b)
}

func countForSeconds(sleepTime time.Duration) {
	endTime := time.Now().Add(sleepTime)
	count := 0
	stringSlice := []string{}
	for time.Now().Before(endTime) {
		randomString := generateRandomString()
		stringSlice = append(stringSlice, randomString)
		count++
		time.Sleep(10 * time.Microsecond)
	}
}

func WithSleep(next http.HandlerFunc, sleepTime time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		countForSeconds(sleepTime)
		next.ServeHTTP(w, req)
	}
}

func WithLog(next http.HandlerFunc) http.HandlerFunc {
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