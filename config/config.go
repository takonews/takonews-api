package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

var Config = struct {
	// debug or release
	Mode string
	DB   struct {
		Name     string
		Adapter  string
		User     string
		Password string
	}
	Secret struct {
		Users []map[string]string
	}
	PORT string
}{}

func init() {
	// debug or release
	Config.Mode = os.Getenv("GIN_MODE")
	if Config.Mode != "debug" && Config.Mode != "release" {
		Config.Mode = "debug"
	}
	// basicauth
	secretConfigPath := os.ExpandEnv("${GOPATH}/src/github.com/takonews/takonews-api/config/secrets.yml")
	buf, err := ioutil.ReadFile(secretConfigPath)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(buf, &Config.Secret)
	if err != nil {
		panic(err)
	}
	// DB
	dbConfigPath := os.ExpandEnv("${GOPATH}/src/github.com/takonews/takonews-api/config/database.yml")
	buf, err = ioutil.ReadFile(dbConfigPath)
	if err != nil {
		panic(err)
	}
	var m interface{}
	err = yaml.Unmarshal(buf, &m)
	if err != nil {
		panic(err)
	}

	if Config.Mode == "debug" {
		development := m.(map[interface{}]interface{})["development"].(map[interface{}]interface{})
		if val, ok := development["name"]; ok {
			Config.DB.Name = val.(string)
		} else {
			Config.DB.Name = ""
		}
		if val, ok := development["adapter"]; ok {
			Config.DB.Adapter = val.(string)
		} else {
			Config.DB.Adapter = ""
		}
		Config.DB.User = development["user"].(string)
		if val, ok := development["user"]; ok {
			Config.DB.User = val.(string)
		} else {
			Config.DB.User = ""
		}
		if val, ok := development["password"]; ok {
			Config.DB.Password = val.(string)
		} else {
			Config.DB.Password = ""
		}
	} else {
		production := m.(map[interface{}]interface{})["production"].(map[interface{}]interface{})
		if val, ok := production["name"]; ok {
			Config.DB.Name = val.(string)
		} else {
			Config.DB.Name = ""
		}
		if val, ok := production["adapter"]; ok {
			Config.DB.Adapter = val.(string)
		} else {
			Config.DB.Adapter = ""
		}
		if val, ok := production["user"]; ok {
			Config.DB.User = val.(string)
		} else {
			Config.DB.User = ""
		}
		if val, ok := production["password"]; ok {
			Config.DB.Password = val.(string)
		} else {
			Config.DB.Password = ""
		}
	}
}
