package config

import (
	"encoding/json"
	"io/ioutil"
	"os/exec"

	"github.com/juju/errors"
)

type Config struct {
	VolumeName string `json:"volume"`
	DeployName string `json:"deploy_name"`

	BaseURL   string `json:"base_url"`
	Port      uint16 `json:"port"`
	HTTPSPort uint16 `json:"https_port"`

	MySQLHost              string `json:"mysql_host"`
	MySQLPort              string `json:"mysql_port"`
	MySQLUserName          string `json:"mysql_username"`
	MySQLPassword          string `json:"mysql_password"`
	ProjectDBName          string `json:"project_db_name"`
	WikiDBName             string `json:"wiki_db_name"`
	IsUseIndependenceMySQL bool   `json:"is_use_independence_mysql"`
	IsUseProxyConfig       bool   `json:"is_use_proxy_config"`
}

func LoadConfig(cfgFile string) (cfg *Config, err error) {
	var b []byte
	b, err = ioutil.ReadFile(cfgFile)
	if err != nil {
		err = errors.Trace(err)
		return
	}
	if err = json.Unmarshal(b, &cfg); err != nil {
		return
	}

	if cfg.MySQLHost == "" {
		cfg.MySQLHost = "localhost"
	}
	if cfg.MySQLPort == "" {
		cfg.MySQLPort = "3306"
	}
	if cfg.ProjectDBName == "" {
		cfg.ProjectDBName = "project"
	}
	if cfg.IsUseProxyConfig {
		var password string
		password, err = InitProxyConfigPassword()
		if err != nil {
			return
		}
		cfg.MySQLPassword = password
	}
	return
}

func InitProxyConfigPassword() (string, error) {
	cmd := exec.Command("/bin/bash", "-c", "./config-proxy/bin/config-proxy -config=./config-proxy/conf/config.json")
	stdout, err := cmd.Output()
	if err != nil {
		return "", errors.Trace(err)
	}
	return string(stdout), nil
}
