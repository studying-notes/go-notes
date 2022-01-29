/*
 * @Date: 2022.01.11 16:01
 * @Description: Omit
 * @LastEditors: Rustle Karl
 * @LastEditTime: 2022.01.11 16:01
 */

package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type HttpConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	BasePath string `mapstructure:"base_path"`
}

type Config struct {
	Http HttpConfig `mapstructure:"http"`

	Logger struct {
		Filename   string `mapstructure:"filename"`
		MaxSize    int    `mapstructure:"max_size"`
		MaxAge     int    `mapstructure:"max_age"`
		MaxBackups int    `mapstructure:"max_backups"`
	} `mapstructure:"logger"`
}

func main() {
	conf := Config{}

	v := viper.New()

	v.SetConfigFile("D:\\OneDrive\\Repositories\\notes\\content\\posts\\go\\src\\libraries\\tripartite\\viper\\config.yaml")

	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := v.Unmarshal(&conf); err != nil {
		panic(err)
	}

	v.OnConfigChange(func(e fsnotify.Event) {
		if err := v.Unmarshal(&conf); err != nil {
			panic(err)
		}
		fmt.Printf("Config file changed: %+v\n", conf)
	})

	v.WatchConfig()

	select {}
}
