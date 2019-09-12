package configs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	KeyWords KeyWordsConfig `json:"keyword"`
	Cookie   string `json:"cookies"`
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
