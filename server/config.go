package server

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"zhibo/kafka"
	"zhibo/mysql"
)

type ApiConfig struct {
	ApiAddress string `yaml:"apiAddress"`
	CookiePath string `yaml:"cookiePath"`
}

type Config struct {
	MysqlConfig *mysql.Config `yaml:"mysql"`
	KafkaConfig *kafka.Config `yaml:"kafka"`
	ApiConfig   *ApiConfig    `yaml:"apiConfig"`
	Api         string        `yaml:"api"`
	QqGroupId   string        `yaml:"qqGroupId"`
}

func GetConfig(path string) *Config {
	//应该是 绝对地址
	yamlFile, err := ioutil.ReadFile(path)
	// 判断是否读取成功
	if err != nil {
		log.Panic(err.Error())
	}
	// 实例化
	config := &Config{}
	// 解码
	err = yaml.Unmarshal(yamlFile, config)
	// 判断解码结果
	if err != nil {
		log.Panic(err.Error())
	}
	// 返回
	return config
}
