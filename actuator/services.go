package actuator

import (
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shumybest/ragnaros/config"
	"github.com/shumybest/ragnaros/eureka"
	"github.com/shumybest/ragnaros/log"
	"github.com/shumybest/ragnaros/repository"
	"github.com/shumybest/ragnaros/utils"
	"net/http"
	"os"
	"strings"
	"time"
)

var logger = log.GetLoggerInstance()

type health struct {
	dbDetails    `json:"db,omitempty"`
	redisDetails `json:"redis,omitempty"`
}

type info struct {
	diskDetails  `json:"diskSpace,omitempty"`
	cpuDetails   `json:"cpuUsage,omitempty"`
	memDetails   `json:"memUsage,omitempty"`
	dbDetails    `json:"db,omitempty"`
	redisDetails `json:"redis,omitempty"`
}

func getDiskSpace() *disk.UsageStat {
	parts, err := disk.Partitions(true)
	if err != nil {
		logger.Errorf("get Partitions failed, err:%v\n", err)
		return nil
	}

	for _, part := range parts {
		if part.Mountpoint == "/" {
			diskInfo, _ := disk.Usage(part.Mountpoint)
			return diskInfo
		}
	}

	return nil
}

func getCpuUsage() []float64 {
	percent, _ := cpu.Percent(time.Second, false)
	return percent
}

func getMemUsage() *mem.VirtualMemoryStat {
	v, _ := mem.VirtualMemory()
	return v
}

func healthHandler(c *gin.Context) {
	dbDetails := dbDetails{}
	if m := repository.GetMySQLInstance(); m.DB != nil {
		dbDetails.Status = m.Status
		dbDetails.Connection = utils.JdbcToDSN(config.GetConfigString("spring.datasource.url"))
	}

	redisDetails := redisDetails{}
	redisPool := repository.GetRedisInstance().Redis.PoolStats()
	if redisPool.TotalConns > 0 {
		redisDetails.Status = eureka.UP
	} else {
		redisDetails.Status = eureka.DOWN
	}
	redisDetails.PoolStatus = *redisPool

	c.JSON(http.StatusOK, gin.H{
		"status": eureka.GetClientInstance().Status,
		"details": health{
			dbDetails,
			redisDetails,
		},
	})
}

func infoHandler(c *gin.Context) {
	// disk details
	diskUsage := getDiskSpace()
	spaceDetail := diskSpace{
		diskUsage.Total,
		diskUsage.Free,
	}

	diskStatus := eureka.UP
	if (float64(diskUsage.Free) / float64(diskUsage.Total)) < 0.2 {
		diskStatus = eureka.DOWN
	}

	diskDetails := diskDetails{
		diskStatus,
		spaceDetail,
	}

	// cpu details
	cpuUsage := getCpuUsage()
	cpuStatus := eureka.UP
	if utils.Average(cpuUsage) > 80.0 {
		cpuStatus = eureka.DOWN
	}
	cpuDetails := cpuDetails{
		cpuStatus,
		cpuUsage,
	}

	// mem details
	vmemUsage := getMemUsage()
	memStatus := eureka.UP
	if vmemUsage.UsedPercent > 80.0 {
		memStatus = eureka.DOWN
	}
	memUsage := memUsage{
		vmemUsage.Total,
		vmemUsage.Free,
		vmemUsage.UsedPercent,
	}
	memDetails := memDetails{
		memStatus,
		memUsage,
	}

	dbDetails := dbDetails{}
	if m := repository.GetMySQLInstance(); m.DB != nil {
		dbDetails.Status = m.Status
		dbDetails.Connection = utils.JdbcToDSN(config.GetConfigString("spring.datasource.url"))
	}

	redisDetails := redisDetails{}
	redisPool := repository.GetRedisInstance().Redis.PoolStats()
	if redisPool.TotalConns > 0 {
		redisDetails.Status = eureka.UP
	} else {
		redisDetails.Status = eureka.DOWN
	}
	redisDetails.PoolStatus = *redisPool

	c.JSON(http.StatusOK, gin.H{
		"status": eureka.GetClientInstance().Status,
		"details": info {
			diskDetails,
			cpuDetails,
			memDetails,
			dbDetails,
			redisDetails,
		},
		// TODO：details when authorized here
	})
}

func envHandler(c *gin.Context) {
	envs := gin.H{}
	for _, v := range os.Environ() {
		vl := strings.Split(v, "=")
		envs[vl[0]] = vl[1]
	}
	c.JSON(http.StatusOK, envs)
}

func configHandler(c *gin.Context) {
	c.JSON(http.StatusOK, config.Context.ConfigStore)
}
