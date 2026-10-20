package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	transmission "github.com/spnngl/transmission-exporter"
)

// RPC responses by method. Torrent 2 lacks file-count, as on Transmission 3.
var rpcResponses = map[string]string{
	"session-get": `{"result":"success","arguments":{"version":"4.0.6","download-dir":"/downloads","download-dir-free-space":1024}}`,
	"session-stats": `{"result":"success","arguments":{"torrentCount":2,
		"current-stats":{"downloadedBytes":10},"cumulative-stats":{"downloadedBytes":20}}}`,
	"torrent-get": `{"result":"success","arguments":{"torrents":[
		{"id":1,"name":"a","file-count":3,"totalSize":5000000000,"trackerStats":[
			{"host":"t.example:443","downloadCount":2,"leecherCount":1,"seederCount":4},
			{"host":"t.example:443","downloadCount":5,"leecherCount":1,"seederCount":4}]},
		{"id":2,"name":"b"}]}}`,
}

func fakeTransmission(t *testing.T) *transmission.Client {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Transmission-Session-Id") != "current" {
			w.Header().Set("X-Transmission-Session-Id", "current")
			w.WriteHeader(http.StatusConflict)
			return
		}
		var req struct {
			Method string `json:"method"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("decode request: %v", err)
		}
		_, _ = w.Write([]byte(rpcResponses[req.Method]))
	}))
	t.Cleanup(srv.Close)
	return transmission.New(srv.URL, nil)
}

// A collector that can't reach Transmission fails the gather, so promhttp
// answers 500 and Prometheus marks the target down.
func TestCollectorsTransmissionDown(t *testing.T) {
	srv := httptest.NewServer(http.NotFoundHandler())
	srv.Close()
	client := transmission.New(srv.URL, nil)

	for name, c := range map[string]prometheus.Collector{
		"torrent":       NewTorrentCollector(client),
		"session":       NewSessionCollector(client),
		"session_stats": NewSessionStatsCollector(client),
	} {
		reg := prometheus.NewPedanticRegistry()
		reg.MustRegister(c)
		if _, err := reg.Gather(); err == nil {
			t.Errorf("%s: gather succeeded with Transmission down", name)
		}
	}
}

// The pedantic registry fails on a metric whose descriptor Describe doesn't send.
func TestCollectors(t *testing.T) {
	client := fakeTransmission(t)
	reg := prometheus.NewPedanticRegistry()
	reg.MustRegister(NewTorrentCollector(client), NewSessionCollector(client), NewSessionStatsCollector(client))

	want := `
# HELP transmission_torrent_files_total The total number of files in a torrent
# TYPE transmission_torrent_files_total gauge
transmission_torrent_files_total{id="1",name="a"} 3
# HELP transmission_torrent_downloads_total How often this torrent was downloaded
# TYPE transmission_torrent_downloads_total gauge
transmission_torrent_downloads_total{id="1",name="a",tracker="t.example:443"} 7
# HELP transmission_torrent_total_size The total size of the torrent
# TYPE transmission_torrent_total_size gauge
transmission_torrent_total_size{id="1",name="a"} 5e+09
transmission_torrent_total_size{id="2",name="b"} 0
# HELP transmission_session_stats_downloaded_bytes The number of downloaded bytes
# TYPE transmission_session_stats_downloaded_bytes gauge
transmission_session_stats_downloaded_bytes{type="cumulative"} 20
transmission_session_stats_downloaded_bytes{type="current"} 10
# HELP transmission_free_space Free space left on device to download to
# TYPE transmission_free_space gauge
transmission_free_space{download_dir="/downloads",incomplete_dir=""} 1024
`
	if err := testutil.GatherAndCompare(reg, strings.NewReader(want),
		"transmission_torrent_files_total",
		"transmission_torrent_downloads_total",
		"transmission_torrent_total_size",
		"transmission_session_stats_downloaded_bytes",
		"transmission_free_space",
	); err != nil {
		t.Fatal(err)
	}
}
