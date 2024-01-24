package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Env         string
	Server      Server
	Log         Log
	DBConfig    DBConfig
	Kafka       Kafka
	RedisConfig RedisConfig
	HTTP        HTTP
	Producer    Producer
	Secret      Secret
	AWSConfig   AWSConfig
	Exception   Exception
}

type Secret struct {
	Private string
}

type Server struct {
	Name string
	Port string
}

type Log struct {
	Level string
}

type DBConfig struct {
	Host            string
	Port            string
	Username        string
	Password        string
	Name            string
	MaxOpenConn     int32
	MaxConnLifeTime int64
}

type Kafka struct {
	Internal KafkaConfig
	External KafkaConfig
}

type RedisConfig struct {
	Mode     string
	Host     string
	Port     string
	Password string
	DB       int
	Cluster  struct {
		Password string
		Addr     []string
	}
}

type KafkaConfig struct {
	Brokers  []string
	Group    string
	Topic    []string
	Producer struct {
		Topic string
	}
	Version  string
	Oldest   bool
	SSAL     bool
	TLS      bool
	CertPath string
	Certs    string
	Username string
	Password string
	Strategy string
}

type HTTP struct {
	TimeOut            time.Duration
	MaxIdleConn        int
	MaxIdleConnPerHost int
	MaxConnPerHost     int
}
type AWSConfig struct {
	RDSSecret    string
	CommonSecret string
	Region       string
}

type Producer struct {
	PaymentListener  string
	FinalTxnListener string
}

type Exception struct {
	Code        ExceptionConfiguration
	Description ExceptionConfiguration
}

type ExceptionConfiguration struct {
	InsufficientGoldBalance string
}

func InitConfig() (*Config, error) {

	configPath, ok := os.LookupEnv("API_CONFIG_PATH")
	if !ok {
		configPath = "./config"
	}

	configName, ok := os.LookupEnv("API_CONFIG_NAME")
	if !ok {
		configName = "config"
	}

	viper.SetConfigName(configName)
	viper.AddConfigPath(configPath)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("config file not found. using default/env config: " + err.Error())
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	var c Config

	err := viper.Unmarshal(&c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func InitTimeZone() {
	ict, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		panic(err)
	}
	time.Local = ict
}
