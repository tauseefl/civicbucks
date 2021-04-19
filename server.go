package main

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"sync"
	"syscall"
	"time"
)

var mutex = &sync.Mutex{}

func echoString(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `Welcome to Civic bucks. 
	Use GET /civicbucks to get civic bucks mined
	Use GET /metrics to get Prometheus formatted metrics
	Use POST /shutdown to shut down mining operation`)
}

func civicCounter(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	var result = getMiningResults()

	var parseDataS struct {
		Total       int
		Palindromes []int
	}
	parseDataS.Total = len(result)

	for i := 0; i < len(result); i++ {
		parseDataS.Palindromes = append(parseDataS.Palindromes, result[i].number)
	}

	b, _ := json.Marshal(parseDataS)
	fmt.Fprintf(w, "%v", string(b))

	mutex.Unlock()
}

func shutdownServer(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request
	// return the string response containing the request body

	switch r.Method {
	case "GET":
		// TODO: Add proper running time
		fmt.Fprintf(w, "Server running since ....")
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		shutdown := r.FormValue("server_shutdown")
		fmt.Fprintf(w, "shutdown = %s\n", shutdown)
		if shutdown == "NOW" {
			log.Print("Force Server Stop received")
			_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		}
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}

}

func getTimeofLastBucks() float64 {
	mutex.Lock()
	var result = getMiningResults()
	mutex.Unlock()
	if len(result) > 1 {
		return float64(result[len(result)-1].time)
	}
	return 0.0
}

var (
	histogram = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "golang",
			Name:      "my_histogram",
			Help:      "This is my histogram",
		})
)

func civicServer() {

	// Create default route handler
	http.HandleFunc("/", echoString)

	http.HandleFunc("/civicbucks", civicCounter)

	http.HandleFunc("/shutdown", shutdownServer)

	http.Handle("/metrics", promhttp.Handler())

	//Registering the defined metric with Prometheus
	//TODO: Complete with actual numbers
	prometheus.MustRegister(histogram)

	go func() {
		for {
			histogram.Observe(getTimeofLastBucks())
			time.Sleep(time.Second)
		}
	}()

	go func() {
		log.Fatal(http.ListenAndServe(":8081", nil))
	}()
}
