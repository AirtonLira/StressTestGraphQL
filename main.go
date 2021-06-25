package main

import (
	"bytes"
	"encoding/json"
	"github.com/StressTestGraphQL/internal/container"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type Config struct {
	PathCertification string `mapstructure:"PATHCERTIFICATION"`
	PathKey			  string `mapstructure:"PATHCERTIFICATIONKEY"`
	Host           string `mapstructure:"HOST"`
}


func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("Query: %s took %s", name,elapsed)
}

func prepareRequest(containerInjector container.Dependency) (client http.Client, host string, queries []string) {

	components := containerInjector.Components
	queries = components.Queries.Query

	transport := components.Transport
	host = components.Conf.Host

	client = http.Client{Transport: transport}

	return client, host, queries


}

func sendRequest(container container.Dependency,queries string, host string, client http.Client, done chan<- bool, wg *sync.WaitGroup) (okay bool){

		defer func(){
			done <- okay
		}()
		body := map[string]string{
			"query":queries}
		jsonValue, _ := json.Marshal(body)

		thumb := container.Components.Conf.Thumbprint

		timeshare := time.Now()
		log.Println("[INFO] - Starting request.... ")
		req, err := http.NewRequest("POST", host, bytes.NewBuffer(jsonValue))
		if err != nil{
			log.Fatalln("[INFO] - SendRequest - Error cert", err)
		}

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("issuer_id","192")
		req.Header.Add("Thumbprint", thumb)

	    resp, err := client.Do(req)
		if err != nil {
			log.Fatalln("Error to send request: ", err)
		}
		log.Println("[INFO] - SendRequest - Status: ", resp.Status, " StatusCode:", resp.StatusCode)
		result, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal("Failed to read req.Body: ",err)
		}


		timeTrack(timeshare, string(result)[1:45])

		req.Body.Close()

		wg.Done()

		return true
}


func main(){
	log.Println("[INFO] - Start program..")

	containerInjector := container.Injector()

	config := containerInjector.Components.Conf

	routines, _  := strconv.ParseInt(config.Goroutines, 10, 32)
	threads, _ := strconv.ParseInt(config.Threads, 10, 32)

	log.Println("[INFO] - Environment PathKey: ",config.PathKey)
	log.Println("[INFO] - Environment PathCertification: ", config.PathCertification)
	log.Println("[INFO] - Environment Limit query: ", config.Limits)
	log.Println("[INFO] - Environment Routines: ", routines)
	log.Println("[INFO] - Environment Threads: ", threads)

	client, host, queries := prepareRequest(containerInjector)

	runtime.GOMAXPROCS(int(threads))

	done := make(chan bool)

	var wg sync.WaitGroup

	for i:= 0; i < int(routines); i++{
		wg.Add(1)
		go func() {
			for _, context := range queries {
				go sendRequest(containerInjector,context, host, client, done, &wg)
			}
		}()
		time.Sleep(1 * time.Second)
	}
	time.Sleep(30 * time.Minute)
}
