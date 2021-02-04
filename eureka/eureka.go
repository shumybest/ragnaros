package eureka

import (
	"bytes"
	"encoding/xml"
	"github.com/shumybest/ragnaros/config"
	"github.com/shumybest/ragnaros/feign"
	. "github.com/shumybest/ragnaros/logger"
	"github.com/shumybest/ragnaros/utils"
	"strings"
	"sync"
	"time"
)

const (
	AppsUrl = "apps/"
)

type Client struct {
	Instance         InstanceConfig
	Status           string
}

var instance *Client
var once sync.Once
func GetClientInstance() *Client {
	once.Do(func() {
		instance = &Client{}
	})
	return instance
}

var eurekaServiceUrl string

func (c *Client) Register() {
	eurekaServiceUrl = config.GetConfigString("eureka.client.service-url.defaultZone")
	if eurekaServiceUrl == "" {
		Logger.Warn("Eureka Service URL is empty, running into mono mode")
		c.Status = OUT_OF_SERVICE
		return
	}

	c.Instance = composeInstance()
	buf, _ := xml.Marshal(c.Instance)
	registerUrl := eurekaServiceUrl + AppsUrl + c.Instance.App

	Logger.Info("trying to register to Eureka: " + registerUrl)

	resp, err := utils.RetryableClient().
		SetHeader("Content-Type", "application/xml").
		SetBody(bytes.NewBuffer(buf)).
		Post(registerUrl)

	if err != nil {
		Logger.Error(err)
		c.Status = OUT_OF_SERVICE
		return
	}

	if resp.StatusCode() == 204 || resp.StatusCode() == 200 {
		Logger.Info("Eureka Client Register Succeed")
		c.Status = UP
		go c.clientRefresh()
	} else {
		c.Status = OUT_OF_SERVICE
		Logger.Warnf("Eureka Client Register Failed: %s %s", resp.StatusCode(), resp)
	}
}

func (c *Client) clientRefresh() {
	defer c.unRegister()
	appsUrl := eurekaServiceUrl + AppsUrl
	instanceUrl := appsUrl + c.Instance.App + "/" + c.Instance.InstanceId

	for {
		// heartbeat
		client := utils.RetryableClient()
		if resp, err := client.Put(instanceUrl); err == nil {
			if resp.StatusCode() != 204 && resp.StatusCode() != 200 {
				c.Status = UNKNOWN
				Logger.Warnf("Eureka Client Renew Failed: %s %s\n", resp)

				// perform register again
				c.unRegister()
				c.Register()
				break
			}
		} else {
			Logger.Error(err)
			c.Status = UNKNOWN
			return
		}

		// get apps
		if resp, err := client.Get(appsUrl); err == nil {
			if resp.StatusCode() == 200 {
				var apps ApplicationsResponse
				if err = xml.Unmarshal(resp.Body(), &apps); err == nil {
					for _, app := range apps.Applications {
						var instances []feign.Instance
						for _, inst := range app.Instances {
							instances = append(instances, feign.Instance{
								HomePageUrls: inst.HomePageUrl,
								Status:       inst.Status,
							})
						}
						feign.Applications[strings.ToLower(app.Name)] = instances
					}
				}
			}
		} else {
			Logger.Error(err)
			c.Status = UNKNOWN
			return
		}

		Logger.Debugf("application refresh: %v\n", feign.Applications)
		config.SetConfig("ragnaros.conf.applications", feign.Applications)
		c.Status = UP
		time.Sleep(10 * time.Second)
	}
}

func (c *Client) unRegister() {
	instanceUrl := eurekaServiceUrl + AppsUrl + c.Instance.App + "/" + c.Instance.InstanceId

	resp, err := utils.RetryableClient().Delete(instanceUrl)
	if err != nil {
		c.Status = UNKNOWN
		Logger.Error(err)
		return
	}

	if resp.StatusCode() != 204 && resp.StatusCode() != 200 {
		c.Status = UNKNOWN
		Logger.Warnf("Eureka Client Delete Failed: %s %s", resp.StatusCode(), resp)
	} else {
		c.Status = DOWN
	}
}
