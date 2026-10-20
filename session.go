package transmission

// Session information about the current transmission session. Speed limits
// are float64: Transmission 4.1 sends them with a fraction (100.0).
type Session struct {
	AltSpeedDown          int     `json:"alt-speed-down"`
	AltSpeedEnabled       bool    `json:"alt-speed-enabled"`
	AltSpeedUp            int     `json:"alt-speed-up"`
	CacheSizeMB           int     `json:"cache-size-mb"`
	DownloadDir           string  `json:"download-dir"`
	DownloadDirFreeSpace  int64   `json:"download-dir-free-space"`
	DownloadQueueEnabled  bool    `json:"download-queue-enabled"`
	DownloadQueueSize     int     `json:"download-queue-size"`
	IncompleteDir         string  `json:"incomplete-dir"`
	PeerLimitGlobal       int     `json:"peer-limit-global"`
	PeerLimitPerTorrent   int     `json:"peer-limit-per-torrent"`
	SeedQueueEnabled      bool    `json:"seed-queue-enabled"`
	SeedQueueSize         int     `json:"seed-queue-size"`
	SeedRatioLimit        float64 `json:"seedRatioLimit"`
	SeedRatioLimited      bool    `json:"seedRatioLimited"`
	SpeedLimitDown        float64 `json:"speed-limit-down"`
	SpeedLimitDownEnabled bool    `json:"speed-limit-down-enabled"`
	SpeedLimitUp          float64 `json:"speed-limit-up"`
	SpeedLimitUpEnabled   bool    `json:"speed-limit-up-enabled"`
	Version               string  `json:"version"`
}
