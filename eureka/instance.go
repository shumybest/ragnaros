package eureka

import (
	"encoding/xml"
	"github.com/shumybest/ragnaros2/config"
	"strings"
)

type DataCenterInfo struct {
	Name string `xml:"name"`
}

type MetaData struct {
	XMLName xml.Name `xml:"metadata"`
	Port    string   `xml:"management.port"`
}

type Instance struct {
	InstanceId                    string `xml:"instanceId"`
	HostName                      string `xml:"hostName"`
	App                           string `xml:"app"`
	IpAddr                        string `xml:"ipAddr"`
	Status                        string `xml:"status"`
	Overriddenstatus              string `xml:"overriddenstatus,omitempty"`
	Port                          string `xml:"port"`
	SecurePort                    string `xml:"securePort"`
	CountryId                     string `xml:"countryId,omitempty"`
	HomePageUrl                   string `xml:"homePageUrl"`
	StatusPageUrl                 string `xml:"statusPageUrl"`
	HealthCheckUrl                string `xml:"healthCheckUrl"`
	VipAddress                    string `xml:"vipAddress"`
	SecureVipAddress              string `xml:"secureVipAddress"`
	IsCoordinatingDiscoveryServer string `xml:"isCoordinatingDiscoveryServer,omitempty"`
	LastUpdatedTimestamp          string `xml:"lastUpdatedTimestamp,omitempty"`
	LastDirtyTimestamp            string `xml:"lastDirtyTimestamp,omitempty"`
	ActionType                    string `xml:"actionType,omitempty"`
}

const (
	UP             = "UP"
	DOWN           = "DOWN"
	STARTING       = "STARTING"
	OUT_OF_SERVICE = "OUT_OF_SERVICE"
	UNKNOWN        = "UNKNOWN"
	MyOwn          = "MyOwn"
)

type InstanceConfig struct {
	Instance
	XMLName        xml.Name       `xml:"instance"` // 指定最外层的标签为instance
	DataCenterInfo DataCenterInfo `xml:"dataCenterInfo"`
	MetaDataInfo   MetaData       `xml:"metadata"`
}

type Application struct {
	Name string `xml:"name"`
	Instances []Instance `xml:"instance"`
}

type ApplicationsResponse struct {
	VersionsDelta string `xml:"versions__delta"`
	AppsHashcode  string `xml:"apps__hashcode"`
	Applications  []Application `xml:"application"`
}

func composeInstance() InstanceConfig {
	var instance InstanceConfig
	preferIpAddress := config.GetConfigBool("eureka.instance.prefer-ip-address")

	if preferIpAddress {
		instance.HostName = config.Context.IpAddr
	} else {
		instance.HostName = config.Context.Hostname
	}

	instance.App = config.GetConfigString("eureka.instance.appname")
	instance.InstanceId = instance.App + "-" + instance.HostName
	instance.IpAddr = config.Context.IpAddr
	instance.VipAddress = strings.ToLower(instance.App)
	instance.SecureVipAddress = strings.ToLower(instance.App)
	instance.Status = UP
	instance.Port = config.Context.Port
	instance.SecurePort = config.Context.Port
	instance.HomePageUrl = "http://" + config.Context.IpAddr + ":" + config.Context.Port + "/"
	instance.StatusPageUrl =  "http://" + config.Context.IpAddr + ":" + config.Context.Port +
		config.Context.Management.StatusPageUrl
	instance.HealthCheckUrl =  "http://" + config.Context.IpAddr + ":" + config.Context.Port +
		config.Context.Management.HealthCheckUrl
	instance.DataCenterInfo = DataCenterInfo{
		Name: MyOwn, // hardcoded value
	}

	return instance
}