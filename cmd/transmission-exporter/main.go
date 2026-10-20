// Command transmission-exporter exposes Transmission metrics to Prometheus.
package main

import (
	"log"
	"net/http"
	"time"

	arg "github.com/alexflint/go-arg"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	transmission "github.com/spnngl/transmission-exporter"
)

// version is the release, set by the release build:
// go build -ldflags "-X main.version=<version>".
var version = "dev"

// Config gets its content from env and passes it on to different packages
type Config struct {
	TransmissionAddr     string `arg:"env:TRANSMISSION_ADDR"`
	TransmissionPassword string `arg:"env:TRANSMISSION_PASSWORD"`
	TransmissionUsername string `arg:"env:TRANSMISSION_USERNAME"`
	WebAddr              string `arg:"env:WEB_ADDR"`
	WebPath              string `arg:"env:WEB_PATH"`
}

// Version is what --version prints (go-arg's Versioned).
func (Config) Version() string {
	return "transmission-exporter " + version
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("no .env present")
	}

	c := Config{
		WebPath:          "/metrics",
		WebAddr:          ":19091",
		TransmissionAddr: "http://localhost:9091",
	}

	arg.MustParse(&c)
	log.Println("starting transmission-exporter", version)

	var user *transmission.User
	if c.TransmissionUsername != "" && c.TransmissionPassword != "" {
		user = &transmission.User{
			Username: c.TransmissionUsername,
			Password: c.TransmissionPassword,
		}
	}

	client := transmission.New(c.TransmissionAddr, user)

	prometheus.MustRegister(NewTorrentCollector(client))
	prometheus.MustRegister(NewSessionCollector(client))
	prometheus.MustRegister(NewSessionStatsCollector(client))

	http.Handle(c.WebPath, promhttp.Handler())

	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`<html>
			<head><title>Transmission Exporter</title></head>
			<body>
			<h1>Transmission Exporter</h1>
			<p><a href="` + c.WebPath + `">Metrics</a></p>
			</body>
			</html>`))
	})

	srv := &http.Server{Addr: c.WebAddr, ReadHeaderTimeout: 10 * time.Second}
	log.Fatal(srv.ListenAndServe())
}

func boolToString(b bool) string {
	if b {
		return "1"
	}
	return "0"
}
