package transmission

import (
	"reflect"
	"strings"
)

type (
	// Torrent represents a transmission torrent. torrent-get returns only the
	// fields asked for: torrentFields asks for every json tag here.
	Torrent struct {
		ID                 int           `json:"id"`
		Name               string        `json:"name"`
		Status             TorrentStatus `json:"status"`
		Added              int64         `json:"addedDate"`
		IsFinished         bool          `json:"isFinished"`
		PercentDone        float64       `json:"percentDone"`
		UploadRatio        float64       `json:"uploadRatio"`
		RateDownload       int64         `json:"rateDownload"`
		RateUpload         int64         `json:"rateUpload"`
		PeersConnected     int           `json:"peersConnected"`
		PeersGettingFromUs int           `json:"peersGettingFromUs"`
		TotalSize          int64         `json:"totalSize"`
		UploadedEver       int64         `json:"uploadedEver"`
		TrackerStats       []TrackerStat `json:"trackerStats"`
		// Nil before Transmission 4.0 (RPC version 17), which added it.
		FileCount *int `json:"file-count"`
	}

	// TorrentStatus is what a torrent is doing: tr_torrent_activity in
	// Transmission's libtransmission/transmission.h.
	TorrentStatus int

	// TrackerStat has stats about the torrent's tracker
	TrackerStat struct {
		Host          string `json:"host"`
		DownloadCount int    `json:"downloadCount"`
		LeecherCount  int    `json:"leecherCount"`
		SeederCount   int    `json:"seederCount"`
	}
)

// torrentStatuses are the TorrentStatus names, by value.
var torrentStatuses = [...]string{
	"stopped",
	"check_pending",
	"checking",
	"download_pending",
	"downloading",
	"seed_pending",
	"seeding",
}

// String is the status name, or "unknown" for one added after Transmission 4.
func (s TorrentStatus) String() string {
	if s < 0 || int(s) >= len(torrentStatuses) {
		return "unknown"
	}
	return torrentStatuses[s]
}

var torrentFields = func() []string {
	t := reflect.TypeFor[Torrent]()
	fields := make([]string, t.NumField())
	for i := range fields {
		fields[i], _, _ = strings.Cut(t.Field(i).Tag.Get("json"), ",")
	}
	return fields
}()
