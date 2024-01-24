package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Env       string
	Server    Server
	Log       Log
	DBConfig  DBConfig
	Kafka     Kafka
	AWSConfig AWSConfig
	Secret    Secret
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

func InitConfig() (*Config, error) {

	viper.SetDefault("ENV", "DEV")

	viper.SetDefault("SERVER.NAME", "reconcile-daily")
	viper.SetDefault("LOG.LEVEL", "info")

	viper.SetDefault("KAFKACONFIG.INTERNAL.BROKERS", "b-1.aspmsknonprod.dlm5z5.c3.kafka.ap-southeast-1.amazonaws.com:9096,b-2.aspmsknonprod.dlm5z5.c3.kafka.ap-southeast-1.amazonaws.com:9096,b-3.aspmsknonprod.dlm5z5.c3.kafka.ap-southeast-1.amazonaws.com:9096")
	viper.SetDefault("KAFKACONFIG.INTERNAL.GROUP", "")
	viper.SetDefault("KAFKACONFIG.INTERNAL.VERSION", "2.8.1")
	viper.SetDefault("KAFKACONFIG.INTERNAL.OLDEST", true)
	viper.SetDefault("KAFKACONFIG.INTERNAL.SSAL", false)
	viper.SetDefault("KAFKACONFIG.INTERNAL.TLS", false)
	viper.SetDefault("KAFKACONFIG.INTERNAL.STRATEGY", "roundrobin")
	viper.SetDefault("HTTP.TIMEOUT", "3s")
	viper.SetDefault("HTTP.MAXIDLECONN", "100")
	viper.SetDefault("HTTP.MAXIDLECONNPERHOST", "100")
	viper.SetDefault("HTTP.MAXCONNPERHOST", "100")

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
