package util

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	PathCertification string `mapstructure:"PATHCERTIFICATION"`
	PathKey			  string `mapstructure:"PATHCERTIFICATIONKEY"`
	Host              string `mapstructure:"HOST"`
	Limits		      string `mapstructure:"TEST_LIMITS"`
	Threads           string `mapstructure:"THREADS""`
	Goroutines        string `mapstructure:"GOROUTINES"`
	Thumbprint		  string `mapstructure:"THUMBPRINT"`
}


func LoadConfig(path string) (config Config){
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("cannot read config:", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil{
		log.Fatal("cannot read config viper:", err)
	}
	return
}
