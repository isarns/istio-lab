
package utils

import (
	"log"
	"os"
	"strconv"
	"strings"
)

func FromStringToIntArray(s string) []int {
	var result []int
	for _, v := range strings.Split(s, ",") {
		s, err := strconv.Atoi(v)
		if err != nil {
			log.Panic("cant convert", v, "to string.")
		}
		result = append(result, s)
	}
	return result
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
