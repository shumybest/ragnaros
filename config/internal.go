package config

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/imdario/mergo"
	"github.com/shumybest/ragnaros/utils"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func shallPutValue(env string, path string, field *string, defaultValue string) {
	envValue := os.Getenv(env)
	if envValue == "" {
		if defaultValue == "" {
			// don't setField
			return
		} else {
			envValue = defaultValue
		}
	}

	if path != "" {
		setField(path, envValue, configStore)
	} else {
		envPath := strings.ToLower(env)
		envPath = strings.Replace(envPath, "_", ".", -1)
		setField(envPath, envValue, configStore)
	}

	if field != nil {
		*field = envValue
	}

	return
}

func loadBootstrapConf() {
	loadYMLConfig(BootstrapConfig, &configStore)
	shallPutValue("SPRING_PROFILES_ACTIVE", "spring.profiles.active", &Context.Profiles, "dev")

	var bootstrapProfileConfig InterfaceMap
	profiles := strings.Split(Context.Profiles, ",")
	for _, p := range profiles {
		loadYMLConfig(BootstrapConfig+"-"+p, &bootstrapProfileConfig)
		mergeConfig(&configStore, bootstrapProfileConfig)
	}
}

func loadApplicationConf() {
	var applicationConf InterfaceMap
	loadYMLConfig(ApplicationConfig, &applicationConf)
	mergeConfig(&configStore, applicationConf)

	var applicationProfileConf InterfaceMap
	profiles := strings.Split(Context.Profiles, ",")
	for _, p := range profiles {
		loadYMLConfig(ApplicationConfig+"-"+p, &applicationProfileConf)
		mergeConfig(&configStore, applicationProfileConf)
	}
}

func fetchSpringCloudConf() {
	shallPutValue("SPRING_CLOUD_CONFIG_URI", "spring.cloud.config.uri", nil, "")
	cloudConfigUrl := GetConfigString("spring.cloud.config.uri")
	cloudConfigLabel := GetConfigString("spring.cloud.config.label")
	cloudConfigProfile := GetConfigString("spring.cloud.config.profile")

	if cloudConfigUrl != "" && cloudConfigLabel != "" && cloudConfigProfile != "" {
		resp, err := utils.RetryableClient().
			Get(cloudConfigUrl + "/" + cloudConfigProfile + "/" + cloudConfigLabel)
		if err != nil {
			fmt.Println(err)
			return
		}

		if resp.StatusCode() == 200 {
			var cloudConfig CloudConfig
			_ = json.Unmarshal([]byte(resp.String()), &cloudConfig)
			for _, source := range cloudConfig.PropertySources {
				fmt.Println("load cloud config " + source.Name)
				for k, v := range source.Source {
					setField(k, v, configStore)
				}
			}
		}
	}
}

func finalizeValue() {
	tmpHostname, err := os.Hostname()
	if err != nil {
		Context.Hostname = utils.RandomString(32)
	} else {
		Context.Hostname = tmpHostname
	}

	Context.IpAddr = utils.GetLocalIp()

	Context.Management.BasePath = GetConfigString("management.endpoints.web.base-path")
	Context.Management.StatusPageUrl = GetConfigString("eureka.instance.status-page-url-path")
	Context.Management.HealthCheckUrl = GetConfigString("eureka.instance.health-check-url-path")

	shallPutValue("EUREKA_CLIENT_SERVICE_URL_DEFAULTZONE",
		"eureka.client.service-url.defaultZone", nil, "")
	shallPutValue("SERVER_PORT", "server.port", &Context.Port, "8999")
	shallPutValue("SPRING_DATASOURCE_URL", "spring.datasource.url", nil, "")
	shallPutValue("SPRING_DATASOURCE_USERNAME", "spring.datasource.username",
		nil, "")
	shallPutValue("SPRING_DATASOURCE_PASSWORD", "spring.datasource.password",
		nil, "")
	shallPutValue("SPRING_REDIS_HOST", "spring.redis.host", nil, "127.0.0.1")
	shallPutValue("SPRING_REDIS_PORT", "spring.redis.port", nil, "6379")

	// elasticsearch address for logging
	shallPutValue("RAGNAROS_ELASTICSEARCH_URL", "", nil, "")
	shallPutValue("RAGNAROS_ELASTICSEARCH_HOST", "", nil, "")
	shallPutValue("RAGNAROS_ELASTICSEARCH_PORT", "", nil, "9200")
	shallPutValue("RAGNAROS_ELASTICSEARCH_USERNAME", "", nil, "")
	shallPutValue("RAGNAROS_ELASTICSEARCH_PASSWORD", "", nil, "")

	base64Secret := GetConfigString("jhipster.security.authentication.jwt.base64-secret")
	secret, _ := base64.StdEncoding.DecodeString(base64Secret)
	Context.Security.JwtSecret = secret

	Context.ConfigStore = &configStore
}

func loadYMLConfig(configFile string, config *InterfaceMap) {
	fullPath := Context.ConfDir + "/" + configFile + SuffixConfig
	buffer, err := ioutil.ReadFile(fullPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = yaml.Unmarshal(buffer, config)
	fmt.Println("load configuration from: " + fullPath)
}

func mergeConfig(dst *InterfaceMap, src InterfaceMap) {
	_ = mergo.Merge(dst, src, mergo.WithOverride)
}

func getField(path string, m InterfaceMap) interface{} {
	fields := strings.Split(path, ".")
	if fields[0] == "" {
		return ""
	}

	for key, value := range m {
		switch value.(type) {
		case string, bool, int, float64:
			if key == fields[0] {
				return value
			}
		case InterfaceMap:
			if key == fields[0] {
				return getField(strings.Join(fields[1:], "."), value.(InterfaceMap))
			}
		}
	}

	return nil
}

func replaceVariable(configStr *string) {
	if strings.Contains(*configStr, "$") {
		// replace ${...}
		re := regexp.MustCompile(`\${([a-zA-Z0-9:.\-_${}]+)}`)
		placeholders := re.FindAllStringSubmatch(*configStr, -1)
		for _, p := range placeholders {
			// has default value
			if strings.Contains(p[1], ":") {
				values := strings.Split(p[1], ":")

				// env value first
				envValue := os.Getenv(values[0])
				if envValue != "" {
					found := re.FindString(*configStr)
					if found != "" {
						*configStr = strings.Replace(*configStr, found, envValue, 1)
					}
					continue
				}

				// configStore value second
				valueInConf := getField(values[0], configStore)
				if valueInConf != nil {
					found := re.FindString(*configStr)
					if found != "" {
						*configStr = strings.Replace(*configStr, found, valueInConf.(string), 1)
					}
					continue
				}

				// default value
				found := re.FindString(*configStr)
				if found != "" {
					*configStr = strings.Replace(*configStr, found, values[1], 1)
				}
			} else {
				found := re.FindString(*configStr)
				if found != "" {
					field := getField(p[1], configStore)
					if field != nil {
						*configStr = strings.Replace(*configStr, found, field.(string), 1)
					} else {
						*configStr = strings.Replace(*configStr, found, "", 1)
					}
				}
			}
			replaceVariable(configStr)
		}
	}
}

func setField(path string, value interface{}, m InterfaceMap) {
	fields := strings.Split(path, ".")
	if fields[0] == "" {
		return
	}

	for k, v := range m {
		switch v.(type) {
		case string, bool, int, float64:
			if k == fields[0] {
				m[k] = value
				return
			}
		case InterfaceMap:
			if k == fields[0] {
				setField(strings.Join(fields[1:], "."), value, v.(InterfaceMap))
				return
			}
		}
	}

	// not existing
	if len(fields) > 1 {
		m[fields[0]] = InterfaceMap{}
		setField(strings.Join(fields[1:], "."), value, m[fields[0]].(InterfaceMap))
	} else {
		m[fields[0]] = value
	}
}
