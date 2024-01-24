package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Env        string
	Server     Server
	Log        Log
	DBConfig   DBConfig
	Kafka      Kafka
	AWSConfig  AWSConfig
	Secret     Secret
	SFTPConfig SFTPConfig
}

type SFTPConfig struct {
	Username     string
	Password     string
	Server       string
	Destination  []DestinationConfig
	FileNameTime string
	Directory    string
}

type DestinationConfig struct {
	PrefixFileName string
	Path           string
	Product        string
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

	viper.SetDefault("KAFKA.INTERNAL.BROKERS", "b-1.aspmsknonprod.dlm5z5.c3.kafka.ap-southeast-1.amazonaws.com:9096,b-2.aspmsknonprod.dlm5z5.c3.kafka.ap-southeast-1.amazonaws.com:9096,b-3.aspmsknonprod.dlm5z5.c3.kafka.ap-southeast-1.amazonaws.com:9096")
	viper.SetDefault("KAFKA.INTERNAL.GROUP", "")
	viper.SetDefault("KAFKA.INTERNAL.VERSION", "2.8.1")
	viper.SetDefault("KAFKA.INTERNAL.OLDEST", true)
	viper.SetDefault("KAFKA.INTERNAL.SSAL", false)
	viper.SetDefault("KAFKA.INTERNAL.TLS", false)
	viper.SetDefault("KAFKA.INTERNAL.STRATEGY", "roundrobin")
	viper.SetDefault("AWSCONFIG.RDSSECRET", "AmazonEKS_RDS_Secret")
	viper.SetDefault("AWSCONFIG.COMMONSECRET", "AmazonEKS_secret")
	viper.SetDefault("AWSCONFIG.REGION", "ap-southeast-1")

	viper.SetDefault("DBCONFIG.MAXOPENCONN", "4")
	viper.SetDefault("DBCONFIG.MAXCONNLIFETIME", "300")

	viper.SetDefault("SFTPConfig.Server", "58.137.161.63")
	viper.SetDefault("SFTPConfig.Username", "ARRTUSR001")
	viper.SetDefault("SFTPConfig.Password", "ARRTP@22Uat")
	viper.SetDefault("SFTPConfig.Directory", "ARRT_NBGW_OUTBOUND")

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
