package exporter

type ServerStatItem struct {
    Tag                   bool       `json:"tag"`
    DiscReadFloat         float64    `json:"discRead"`
    DiscWriteFloat        float64    `json:"discWrite"`
    MemoryUsageRate       float64    `json:"memoryUsageRate"`
    DiscInfoList          []DiskInfo `json:"discInfo"`
    CpuUsageRate          float64    `json:"cpuUsageRate"`
    Load                  float64    `json:"load"`
    BandwidthUpload       float64    `json:"bandwidthUpload"`
    BandwidthDownload     float64    `json:"bandwidthDownload"`
    DeviceLinkingCountInt int64      `json:"deviceLinkingCount"`
    IopsRead              float64    `json:"iopsRead"`
    IopsWrite             float64    `json:"iopsWrite"`
}

type DiskInfo struct {
    DiscName  string
    Total     float64
    UsageRate float64
}

type Number struct {
	CpuUsageRate			float64
	MemoryUsageRate			float64
	DiskUsage				float64
	Load					float64
	DiscRead				float64
	DiscWrite				float64
	IopsRead				float64
	IopsWrite				float64	
	DeviceLinkingCount		int64
}
