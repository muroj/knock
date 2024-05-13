package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"
)

type AccessLogEntry struct {
	AccessTime     string `json:"time"`
	RemoteIP       string `json:"remote_ip"`
	RemoteUser     string `json:"remote_user"`
	Request        string `json:"request"`
	HttpStatusCode int    `json:"response"`
	ResponseSize   int    `json:"bytes"`
}

type Metrics struct {
	StatusCodeCounts    map[string]int       `json:"status_code_metrics"`
	Total               *ResponseSizeMetrics `json:"total_response_metrics"`
	Failed              *ResponseSizeMetrics `json:"failed_response_metrics"`
	Successful          *ResponseSizeMetrics `json:"success_response_metrics"`
	MaxResponseSizePath string               `json:"max_response_size_path"`
	BadEndpoint         string               `json:"bad_endpoint"`
}

type ResponseSizeMetrics struct {
	Mean   int `json:"mean"`
	Median int `json:"median"`
	P99    int `json:"p99"`
}

func main() {
	logfilePath := "nginx_json_logs"
	var results *Metrics
	var jsonOut []byte

	f, err := os.Open(logfilePath)
	if err != nil {
		log.Fatalf("failed to open file %s: %v", logfilePath, err)
	}
	defer f.Close()

	results, err = ParseAccessLogs(f)
	if err != nil {
		log.Fatal(err)
	}

	jsonOut, err = json.Marshal(results)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v", string(jsonOut))
}

func ParseAccessLogs(logfile *os.File) (*Metrics, error) {
	var entry AccessLogEntry
	var total, failed, successful []int
	var err error
	var maxResponseSizePath, mostFailuresPath string
	endpointFailureCount := make(map[string]int)
	statusCodeCount := make(map[string]int)
	maxResponseSize := -1
	mostFailures := -1

	dec := json.NewDecoder(logfile)
	for dec.More() {
		err = dec.Decode(&entry)
		if err != nil {
			fmt.Print(err)
			continue
		}

		httpStatus := http.StatusText(entry.HttpStatusCode)
		if httpStatus != "" {
			statusCodeCount[httpStatus]++
		}

		requestFields := strings.Fields(entry.Request)
		httpPath := requestFields[1]

		if entry.HttpStatusCode >= 400 {
			failed = append(failed, entry.ResponseSize)
			endpointFailureCount[httpPath]++
		} else {
			successful = append(successful, entry.ResponseSize)
		}
		total = append(total, entry.ResponseSize)

		if entry.ResponseSize > maxResponseSize {
			maxResponseSize = entry.ResponseSize
			maxResponseSizePath = httpPath
		}
	}

	for k, v := range endpointFailureCount {
		if v > mostFailures {
			mostFailures = v
			mostFailuresPath = k
		}
	}

	return &Metrics{
		Total:               ComputeResults(total),
		Failed:              ComputeResults(failed),
		Successful:          ComputeResults(successful),
		MaxResponseSizePath: maxResponseSizePath,
		BadEndpoint:         mostFailuresPath,
		StatusCodeCounts:    statusCodeCount,
	}, nil
}

func ComputeResults(vals []int) *ResponseSizeMetrics {
	return &ResponseSizeMetrics{
		Mean:   mean(vals),
		Median: median(vals),
		P99:    p99(vals),
	}
}

func mean(vals []int) int {
	sum := 0
	for _, v := range vals {
		sum += v
	}
	return sum / len(vals)
}

func median(vals []int) int {
	slices.Sort(vals)
	middle := (len(vals) / 2) - 1
	return vals[middle]
}

func p99(vals []int) int {
	slices.Sort(vals)
	p99_index := int(float64(len(vals)) * 0.99)
	return vals[p99_index]
}
