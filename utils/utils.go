package utils

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"math/rand"
	"net"
	"net/url"
	"regexp"
	"time"
)

func GetLocalIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		return ""
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	return ""
}

const letterBytes = "abcdefghijklmnopqrstuvwxyz0123456789"

func RandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func Average(s []float64) float64 {
	total := 0.0
	for _, i := range s {
		total += i
	}

	return total / float64(len(s))
}

func JdbcToDSN(connStr string) string {
	// jdbc:mysql://rm-bp1f9x621b96i956d.mysql.rds.aliyuncs.com:3306/datacenter
	// to
	// tcp(rm-bp1f9x621b96i956d.mysql.rds.aliyuncs.com:3306)/datacenter
	// parameters:
	//   useUnicode=true => remove
	//   characterEncoding=utf8 => charset=utf8
	//   useSSL=false => tls=false
	//   useLegacyDatetimeCode=false => remove
	//   serverTimezone=Asia/Shanghai => time_zone=Asia/Shanghai

	ret := ""
	jdbcUrl, _ := url.Parse(connStr)
	re := regexp.MustCompile(`mysql://([a-zA-Z0-9.:\-]+)/(.+)`)
	dest := re.FindAllStringSubmatch(jdbcUrl.Opaque, -1)
	if len(dest) == 1 && len(dest[0]) == 3 {
		ret = "tcp(" + dest[0][1] + ")/" + dest[0][2]
	}

	jdbcQuery, _ := url.ParseQuery(jdbcUrl.RawQuery)
	dsnParams := url.Values{}
	for key, value := range jdbcQuery {
		switch key {
		case "characterEncoding":
			dsnParams.Add("charset", value[0])
		case "useSSL":
			dsnParams.Add("tls", value[0])
		case "serverTimezone":
			dsnParams.Add("time_zone", "'"+value[0]+"'")
			// TODO: more jdbc params conversion to be added
		}
	}

	return ret + "?" + dsnParams.Encode()
}

func RetryableClient() *resty.Request {
	client := resty.New()
	client.
		SetRetryCount(2).
		SetRetryWaitTime(300 * time.Millisecond)

	return client.R()
}

