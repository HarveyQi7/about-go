package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Content struct {
	GinMode        string      `json:"ginMode"`
	Datasources    Datasources `json:"datasources"`
	TrustedProxies []string    `json:"trustedProxies"`
	JwtSecretKey   string      `json:"jwtSecretKey"`
}

var ctn Content

const (
	LOCAL_GO_ENV_MODE = "local"
	DEV_GO_ENV_MODE   = "development"
	PRO_GO_ENV_MODE   = "production"
)

func Load() Content {
	file, err := ioutil.ReadFile("config/" + GetGoEnvMode() + ".json")
	if err != nil {
		panic(err)
	}
	if err = json.Unmarshal(file, &ctn); err != nil {
		panic(err)
	}
	return ctn
}

func Get() Content {
	return ctn
}

func GetGoEnvMode() string {
	value := os.Getenv("go-env-mode")
	if value == "" {
		value = LOCAL_GO_ENV_MODE
	}
	return value
}
