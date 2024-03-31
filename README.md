# Go Status Checker

A VERY simple go app that takes in a new line separated list of URls and reports on the status code & response time, built out of my learning of goroutines, channels etc.

### Docker compose
A docker compose file is used to spin up the container image along with Prometheus & Grafana.

### Running locally
To run the Status checker locally:
```
go run main.go -a </path/to/file> -s <seconds>
```
