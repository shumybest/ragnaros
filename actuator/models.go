package actuator

import "github.com/go-redis/redis/v8"

type diskDetails struct {
	Status    string `json:"status"`
	diskSpace `json:"details"`
}

type diskSpace struct {
	Total uint64 `json:"total"`
	Free  uint64 `json:"free"`
}

type cpuDetails struct {
	Status   string    `json:"status"`
	CpuUsage []float64 `json:"details"`
}

type memDetails struct {
	Status   string `json:"status"`
	memUsage `json:"details"`
}

type memUsage struct {
	Total       uint64  `json:"total"`
	Available   uint64  `json:"available"`
	UsedPercent float64 `json:"percent"`
}

type dbDetails struct {
	Status   string `json:"status"`
	dbStatus `json:"details"`
}

type dbStatus struct {
	Connection string `json:"connection"`
}

type redisDetails struct {
	Status      string `json:"status"`
	redisStatus `json:"details"`
}

type redisStatus struct {
	PoolStatus redis.PoolStats `json:"pool"`
}
