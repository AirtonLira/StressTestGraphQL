package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/StressTestGraphQL/internal/container"
	"github.com/StressTestGraphQL/metric"
	"github.com/google/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const logPath = "./log.log"

var verbose = flag.Bool("verbose", false, "print info level logs to stdout")

type Config struct {
	PathCertification string `mapstructure:"PATHCERTIFICATION"`
	PathKey           string `mapstructure:"PATHCERTIFICATIONKEY"`
	Host              string `mapstructure:"HOST"`
}

func prepareRequest(containerInjector container.Dependency) (client http.Client, host string, queries []string) {

	components := containerInjector.Components
	queries = components.Queries.Query

	transport := components.Transport
	host = components.Conf.Host

	client = http.Client{Transport: transport}

	return client, host, queries

}

func sendRequest(container container.Dependency, queries string, host string, client http.Client, done chan<- bool, wg *sync.WaitGroup, domain string) (okay bool) {
	defer func() {
		done <- okay
	}()
	body := map[string]string{
		"query": queries}
	jsonValue, _ := json.Marshal(body)

	thumb := container.Components.Conf.Thumbprint

	timeshare := time.Now()
	appMetric := metric.HTTP{}

	req, err := http.NewRequest("POST", host, bytes.NewBuffer(jsonValue))
	if err != nil {
		logger.Info("SendRequest - Error cert", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("issuer_id", "192")
	req.Header.Add("Thumbprint", thumb)

	resp, err := client.Do(req)
	if err != nil {
		logger.Fatalf("Error to send request: ", err)
	}

	appMetric.TimeTrack(timeshare, domain, resp.StatusCode)

	req.Body.Close()

	wg.Done()

	return true
}

func main() {

	lf, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		logger.Fatalf("Failed to open log file: %v", err)
	}
	defer lf.Close()

	defer logger.Init("StressTest - GraphQL", *verbose, true, lf).Close()
	logger.Info("Start program..")

	containerInjector := container.Injector()

	config := containerInjector.Components.Conf

	routines, _ := strconv.ParseInt(config.Goroutines, 10, 32)
	threads, _ := strconv.ParseInt(config.Threads, 10, 32)
	hostUrl := config.Host
	domain := config.Domain

	logger.Info("Environment PathKey: ", config.PathKey)
	logger.Info("Environment PathCertification: ", config.PathCertification)
	logger.Info("Environment Limit query: ", config.Limits)
	logger.Info("Environment Routines: ", routines)
	logger.Info("Environment Threads: ", threads)
	logger.Info("Environment URL: ", hostUrl)

	client, host, queries := prepareRequest(containerInjector)

	runtime.GOMAXPROCS(int(threads))

	done := make(chan bool)

	var wg sync.WaitGroup

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8080", nil)
	}()

	logger.Info("Starting request.... ")
	for i := 0; i < int(routines); i++ {
		wg.Add(1)
		go func() {
			for _, context := range queries {
				go sendRequest(containerInjector, context, host, client, done, &wg, domain)
			}
		}()
		time.Sleep(1 * time.Second)
	}
	time.Sleep(30 * time.Minute)
}
