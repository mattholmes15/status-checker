package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var (
	websiteResponseCode = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "website_response_code",
		Help: "The HTTP response code for the website",
	}, []string{"website"})

	websiteResponseTime = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "website_response_time_seconds",
		Help: "The response time for the website in seconds",
	}, []string{"website"})
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	filePath := flag.String("a", "", "Path to the file containing website URLs")
	timeInterval := flag.Int64("s", 0, "Seconds between each check")
	flag.Parse()

	if *filePath == "" {
		fmt.Println("Please provide a file path with the -a flag.")
		return
	}

	file, err := os.Open(*filePath)
	if err != nil {
		fmt.Printf("Error opening file: %s\n", err)
		return
	}
	defer file.Close()

	var websites []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())
		if url != "" {
			websites = append(websites, url)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		return
	}

	c := make(chan string)

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":8090", nil); err != nil {
			fmt.Printf("Error starting server: %s\n", err)
		}
	}()

	for _, w := range websites {
		go checkLink(w, c)
	}

	for w := range c {
		go func(website string) {
			time.Sleep(time.Second * time.Duration(*timeInterval))
			checkLink(website, c)
		}(w)
	}
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func checkLink(website string, c chan string) {
	start := time.Now()
	res, err := http.Get(website)
	duration := time.Since(start).Seconds()
	roundedDuration := roundFloat(duration, 3)

	if err != nil {
		websiteResponseCode.WithLabelValues(website).Set(0)
		websiteResponseTime.WithLabelValues(website).Set(roundedDuration)
 		log.WithFields(
        	log.Fields{
            	"website": website,
            	"status": res.StatusCode,
				"response_time": roundedDuration,
        	},
		).Error(website, "is down!")
		c <- website
		return
	}
	websiteResponseCode.WithLabelValues(website).Set(float64(res.StatusCode))
	websiteResponseTime.WithLabelValues(website).Set(roundedDuration)
    log.WithFields(
        log.Fields{
            "website": website,
            "status": res.StatusCode,
			"response_time": roundedDuration,
        },
	).Info("Success")
	c <- website
}



