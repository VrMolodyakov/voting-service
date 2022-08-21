package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Port       string  `yaml : "port"`
	Host       string  `yaml : "host"`
	LogLvl     string  `yaml : "loglvl"`
	PostgreSql Postgre `yaml : "postgresql"`
	Redis      Redis   `yaml : "redis"`
}

type Redis struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	DbNumber int    `yaml:"dbnumber"`
}

type Postgre struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Dbname   string `yaml:"dbname"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	PoolSize string `yaml:"poolsize"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		path, _ := os.Getwd()
		fmt.Println("path:", path)
		root := filepath.Dir(filepath.Dir(path))
		fmt.Println("dir2:", root)
		instance = &Config{}
		log.Println("start config initialisation")
		configPath := root + "\\config\\config.yaml"
		dockerPath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		if exist, _ := Exists(configPath); exist {
			if err := cleanenv.ReadConfig(root+"\\config\\config.yaml", instance); err != nil {
				log.Fatal(err)
			}
		} else if exist, _ := Exists(dockerPath + "/config/config.yaml"); exist {
			if err := cleanenv.ReadConfig(dockerPath+"/config/config.yaml", instance); err != nil {
				log.Fatal(err)
			}
		}

	})
	return instance
}

func Exists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}
