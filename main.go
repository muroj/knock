package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type AccessLogEntry struct {
	AccessTime       string `json:"time"`
	RemoteIP         string `json:"remote_ip"`
	RemoteUser       string `json:"remote_user"`
	Request          string `json:"request"`
	HttpResponseCode int32  `json:"response"`
	ResponseSize     int64  `json:"bytes"`
}

type Results struct {
}

func main() {
	logfile := "nginx_json_logs"
	f, err := os.Open(logfile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	var entry AccessLogEntry

	for dec.More() {
		err = dec.Decode(&entry)
		if err != nil {
			fmt.Print(err)
			continue
		}

		fmt.Printf("%v", entry.RemoteIP)
	}

}
