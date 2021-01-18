package feign

import (
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/shumybest/ragnaros2/config"
	"strings"
	"time"
)

type Client struct {
	*resty.Client
	App string
}

type Instance struct {
	HomePageUrls string
	Status       string
}
var Applications = make(map[string][]Instance)

type Headers map[string]string
var roundRobin = make(map[string]int)

func getTimeoutConfiguration() time.Duration {
	timeout := config.GetConfigInt("feign.client.config.default.connectTimeout")
	timeout += config.GetConfigInt("feign.client.config.default.readTimeout")

	if timeout == 0 {
		return 5 * time.Second
	}
	return time.Duration(timeout)
}

func getNextServer(app string) string {
	lowerApp := strings.ToLower(app)
	instances := Applications[lowerApp]

	instanceCount := len(instances)
	if instanceCount == 0 {
		return ""
	}

	current := instances[roundRobin[lowerApp]]
	roundRobin[lowerApp]++
	if roundRobin[lowerApp] >= instanceCount {
		roundRobin[lowerApp] = 0
	}
	next := instances[roundRobin[lowerApp]]

	if next.Status == "UP" {
		return next.HomePageUrls
	} else if current.Status == "UP" {
		return current.HomePageUrls
	} else {
		return ""
	}
}

func App(app string) *Client {
	c := Client{App: app}

	c.Client = resty.New()
	c.SetRetryCount(2).SetRetryWaitTime(getTimeoutConfiguration() * time.Millisecond)

	return &c
}

func (c *Client) SetHeaders(headers Headers) {
	c.Client.SetHeaders(headers)
}

// TODO: refine the codes, some kind of ugly
func (c *Client) Get(path string) (*resty.Response, error) {
	dst := getNextServer(c.App)
	if dst == "" {
		return nil, errors.New("server " + c.App + " not available")
	}

	c.SetHostURL(dst)
	resp, err := c.R().Get(path)
	if err == nil {
		return resp, nil
	}

	// try next server on error
	c.SetHostURL(getNextServer(c.App))
	resp, err = c.R().Get(path)
	if err == nil {
		return resp, nil
	}

	return nil, err
}

func (c *Client) AsyncGet(path string, callback func(*resty.Response)) error {
	ch := make(chan *resty.Response, 1)

	dst := getNextServer(c.App)
	if dst == "" {
		return errors.New("server " + c.App + " not available")
	}

	go func() {
		c.SetHostURL(dst)
		resp, err := c.R().Get(path)
		if err == nil {
			ch <- resp
			return
		}

		c.SetHostURL(getNextServer(c.App))
		resp, err = c.R().Get(path)
		if err == nil {
			ch <- resp
		}
		defer close(ch)
	}()

	go func() {
		resp := <- ch
		if callback != nil {
			callback(resp)
		}
		defer close(ch)
	}()

	return nil
}

func (c *Client) Post(path string) (*resty.Response, error) {
	dst := getNextServer(c.App)
	if dst == "" {
		return nil, errors.New("server " + c.App + " not available")
	}

	c.SetHostURL(dst)
	resp, err := c.R().Post(path)
	if err == nil {
		return resp, nil
	}

	// try next server on error
	c.SetHostURL(getNextServer(c.App))
	resp, err = c.R().Post(path)
	if err == nil {
		return resp, nil
	}

	return nil, err
}

func (c *Client) AsyncPost(path string, callback func(*resty.Response)) error {
	ch := make(chan *resty.Response)

	dst := getNextServer(c.App)
	if dst == "" {
		return errors.New("server " + c.App + " not available")
	}

	go func() {
		c.SetHostURL(dst)
		resp, err := c.R().Post(path)
		if err == nil {
			ch <- resp
			return
		}

		c.SetHostURL(getNextServer(c.App))
		resp, err = c.R().Post(path)
		if err == nil {
			ch <- resp
		}
		defer close(ch)
	}()

	go func() {
		resp := <- ch
		if callback != nil {
			callback(resp)
		}
		defer close(ch)
	}()

	return nil
}
