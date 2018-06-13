package models

import (
	"io/ioutil"
	"github.com/BurntSushi/toml"
	"log"
)

type Config struct {
	DB DBConfig `toml: "dbs"`
}

type DBConfig struct {
	Name         string   `toml: "name"`
	Account      string   `toml: "account"`
	Port         int      `toml: "port"`
	Ips          []string `toml: "ips"`
	MaxIdleConns int      `toml: "maxIdleConns"`
	MaxOpenConns int      `toml: "maxOpenConns"`
}

var config * Config

func Conf() *Config {
	return config
}

func InitConfig() {

	if bs, err := ioutil.ReadFile("./conf/db.cfg"); err != nil {
		log.Fatalf("read config file failed: %s\n", err.Error())
	} else {
		if _, err := toml.Decode(string(bs), &config); err != nil {
			log.Fatalf("decode config file failed: %s\n", err.Error())
		} else {
			log.Printf("load config from ./dbs.cfg :\n%v\n", config)
		}
	}
	log.Printf("配置文件内容：%#v", config)
}