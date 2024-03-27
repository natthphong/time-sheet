package config

import (
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Env         string
	EnableS3    bool
	Server      Server
	Log         Log
	DBConfig    DBConfig
	Kafka       Kafka
	AWSConfig   AWSConfig
	Secret      Secret
	RedisConfig RedisConfig
	S3Config    S3Config
	Producer    Producer
}

type Producer struct {
	GoldListener     string
	FinalTxnListener string
}
type S3Config struct {
	BucketName string
	Key        string
}

type ToggleConfiguration struct {
	IsTest bool
	Case   string
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
	CertFile           []byte
	KeyFile            []byte
}
type AWSConfig struct {
	RDSSecret    string
	CommonSecret string
	Region       string
}

func InitConfig() (*Config, error) {

	viper.SetDefault("ENV", os.Getenv("ENV"))
	viper.SetDefault("EnableS3", os.Getenv("enableS3"))
	viper.SetDefault("SERVER.NAME", "his_pricing")
	//viper.SetDefault("LOG.LEVEL", os.Getenv("logLevel"))
	viper.SetDefault("LOG.LEVEL", "debug")

	viper.SetDefault("KAFKA.INTERNAL.BROKERS", os.Getenv("kafkaBrokers"))
	viper.SetDefault("KAFKA.INTERNAL.GROUP", "")
	viper.SetDefault("KAFKA.INTERNAL.VERSION", "2.8.1")
	viper.SetDefault("KAFKA.INTERNAL.OLDEST", true)
	viper.SetDefault("KAFKA.INTERNAL.SSAL", true)
	viper.SetDefault("KAFKA.INTERNAL.TLS", true)
	viper.SetDefault("KAFKA.INTERNAL.STRATEGY", "roundrobin")
	viper.SetDefault("AWSCONFIG.RDSSECRET", os.Getenv("awsRdsSecret"))
	viper.SetDefault("AWSCONFIG.COMMONSECRET", os.Getenv("awsCommonSecret"))
	viper.SetDefault("AWSCONFIG.REGION", "ap-southeast-1")

	viper.SetDefault("DBCONFIG.MAXOPENCONN", "4")
	viper.SetDefault("DBCONFIG.MAXCONNLIFETIME", "300")
	viper.SetDefault("DBCONFIG.Host", "aurora-nonprod-iam-db.cberwwykerv8.ap-southeast-1.rds.amazonaws.com")
	viper.SetDefault("DBCONFIG.Name", "iam_db")
	viper.SetDefault("DBCONFIG.Port", "5432")
	viper.SetDefault("DBCONFIG.Username", "ibm_app")
	viper.SetDefault("DBCONFIG.Password", "[Q]sb3pl*7r*xa7]")

	//viper.SetDefault("S3Config.BucketName", os.Getenv("bucketName"))
	//viper.SetDefault("S3Config.Key", os.Getenv("keyBucket"))

	viper.SetDefault("S3Config.BucketName", "poc-sync-app")
	viper.SetDefault("S3Config.Key", "his_pricing_%s.zip")

	//viper.SetDefault("REDISCONFIG.MODE", os.Getenv("redisMode"))
	//viper.SetDefault("REDISCONFIG.HOST", os.Getenv("redisHost"))
	//viper.SetDefault("REDISCONFIG.Cluster.Addr", os.Getenv("redisHost"))
	//viper.SetDefault("REDISCONFIG.PORT", "6379")

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
