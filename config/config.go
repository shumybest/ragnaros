package config

import (
	"github.com/dgrijalva/jwt-go"
)

// tricky interface wrapper
type InterfaceMap map[string]interface{}
var configStore = make(InterfaceMap)

// spring cloud config
type propertySource struct {
	Name string `json:"name"`
	Source map[string]interface{} `json:"source"`
}
type CloudConfig struct {
	PropertySources []propertySource `json:"propertySources"`
}

// config cache
type Management struct {
	BasePath       string
	HealthCheckUrl string
	StatusPageUrl  string
}

type Security struct {
	JwtSecret []byte
	Claims
}

type Claims struct {
	Auth string `json:"auth"`
	Tid uint `json:"tid"`
	Uid uint `json:"uid"`
	Rid string `json:"rid"`
	jwt.StandardClaims
}

var Context = struct {
	Profiles           string
	ConfDir            string
	Hostname           string
	IpAddr             string
	Port               string

	Management
	Security

	// configuration cache
	ConfigStore *InterfaceMap
}{}

const (
	BootstrapConfig   = "bootstrap"
	ApplicationConfig = "application"
	SuffixConfig      = ".yml"
)

func InitConfig() {
	shallPutValue("RAGNAROS_CONF_DIR", "", &Context.ConfDir, "resources/config")

	loadBootstrapConf()
	loadApplicationConf()
	fetchSpringCloudConf()
	finalizeValue()
}

func GetConfigInt(path string) int {
	field := getField(path, configStore)
	if field != nil {
		return field.(int)
	}
	return 0
}

func GetConfigString(path string) string {
	field := getField(path, configStore)
	if field != nil {
		configStr := field.(string)
		replaceVariable(&configStr)
		return configStr
	}

	return ""
}

func GetConfigBool(path string) bool {
	field := getField(path, configStore)
	if field != nil {
		configValue := field.(bool)
		return configValue
	}

	return false
}

func SetConfig(path string, data interface{}) {
	setField(path, data, configStore)
}
