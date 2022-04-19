package bot

import (
	"io/ioutil"
	"log"
)
import "gopkg.in/yaml.v2"

type Config struct {
	KafkaAddress string `yaml:"address"`
	Topic        string `yaml:"topic"`
	GroupId      string `yaml:"groupId"`
	Api          string `yaml:"api"`
	QqGroupId    string `yaml:"qqGroupId"`
}

func NewConfigWithFile(file string) (*Config, error) {
	//应该是 绝对地址
	yamlFile, err := ioutil.ReadFile(file)
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
	return config, nil
}
