package configs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"runtime"
)

type DataBaseConfig struct {
	Host   string `json:"host"`
	Port   int    `json:"port"`
	User   string `json:"user"`
	Pass   string `json:"pass"`
	DBname string `json:"dbname"`
}

type KeyWordsConfig struct {
	EastMoney []string `json:"eastmoney"`
	TouTiao   []string `json:"toutiao"`
}

type Config struct {
	DataBase DataBaseConfig `json:"database"`
	KeyWords KeyWordsConfig `json:"keywords"`
	Cookie   string         `json:"cookies"`
}

func (c *Config) ReadConfig() {
	yamlFile, err := ioutil.ReadFile("./conf.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(yamlFile, &c)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func ChromePath() string {
	sys := runtime.GOOS
	var chrome string
	switch sys {
	case "darwin":
		chrome = path.Join("./plugins", "chromedriver_mac")
	case "windows":
		chrome = path.Join("./plugins", "chromedriver_win.exe")
	case "linux":
		chrome = path.Join("./plugin", "chromedriver_linux")
	}
	return chrome
}

func Platform() string {
	sys := runtime.GOOS
	var platform string
	switch sys {
	case "darwin":
		platform = "Mac"
	case "windows":
		platform = "Windows"
	case "linux":
		platform = "Linux"
	}
	return platform
}
