package config

import (
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Env                string
	EnableS3           bool
	Server             Server
	Log                Log
	DBConfig           DBConfig
	Kafka              Kafka
	AWSConfig          AWSConfig
	Secret             Secret
	SFTPConfig         SFTPConfig
	RedisConfig        RedisConfig
	Toggle             Toggle
	FundTransferConfig FundTransferConfig
	Exception          Exception
	HTTP               HTTP
	S3Config           S3Config
	Producer           Producer
}

type Producer struct {
	GoldListener     string
	FinalTxnListener string
}
type S3Config struct {
	BucketName string
	Key        string
}

type Exception struct {
	Code        ExceptionConfiguration
	Description ExceptionConfiguration
}

type ExceptionConfiguration struct {
	SystemError string
}

type Toggle struct {
	OauthFundTransfer         ToggleConfiguration
	InquiryStatusFundTransfer ToggleConfiguration
}
type ToggleConfiguration struct {
	IsTest bool
	Case   string
}
type SFTPConfig struct {
	Username     string
	Password     string
	Server       string
	Destination  DestinationConfig
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
	CertFile           []byte
	KeyFile            []byte
}
type AWSConfig struct {
	RDSSecret    string
	CommonSecret string
	Region       string
}

type FundTransferConfig struct {
	OauthFundTransferUrl         string
	InquiryStatusFundTransferUrl string
	ConsumerID                   string
	ConsumerSecret               string
	MerchantID                   string
	FromAccountNo                string
	Auth                         string
	SenderName                   string
	SenderTaxID                  string
	TypeOfSender                 string
	OAuthLimit                   int64
	OauthRetry                   int
	InquiryStatusRetry           int
}

func InitConfig() (*Config, error) {

	viper.SetDefault("ENV", os.Getenv("ENV"))
	viper.SetDefault("EnableS3", os.Getenv("enableS3"))
	viper.SetDefault("SERVER.NAME", "reconcile-daily")
	viper.SetDefault("LOG.LEVEL", "info")
	viper.SetDefault("EXCEPTION.CODE.SystemError", "999")
	viper.SetDefault("EXCEPTION.DESCRIPTION.SystemError", "System error")

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
	viper.SetDefault("DBCONFIG.Name", "postgres")

	viper.SetDefault("SFTPConfig.Server", os.Getenv("sftpHost"))
	viper.SetDefault("SFTPConfig.Username", os.Getenv("sftpUsername"))
	viper.SetDefault("SFTPConfig.Password", os.Getenv("sftpPassword"))
	viper.SetDefault("SFTPConfig.Directory", os.Getenv("sftpDirectory"))
	viper.SetDefault("SFTPConfig.Destination.PrefixFileName", os.Getenv("sftpPrefixName"))
	viper.SetDefault("SFTPConfig.Path", "/Path")
	viper.SetDefault("SFTPConfig.Product", "Product")
	viper.SetDefault("S3Config.BucketName", os.Getenv("bucketName"))
	viper.SetDefault("S3Config.Key", os.Getenv("keyBucket"))

	viper.SetDefault("REDISCONFIG.MODE", os.Getenv("redisMode"))
	viper.SetDefault("REDISCONFIG.HOST", os.Getenv("redisHost"))
	viper.SetDefault("REDISCONFIG.Cluster.Addr", os.Getenv("redisHost"))
	viper.SetDefault("REDISCONFIG.PORT", "6379")

	viper.SetDefault("TOGGLE.OAUTHFUNDTRANSFER.ISTEST", os.Getenv("authTest"))
	viper.SetDefault("TOGGLE.OAUTHFUNDTRANSFER.CASE", os.Getenv("authTestCase"))
	viper.SetDefault("TOGGLE.INQUIRYSTATUSFUNDTRANSFER.ISTEST", os.Getenv("inqTest"))
	viper.SetDefault("TOGGLE.INQUIRYSTATUSFUNDTRANSFER.CASE", os.Getenv("inqTestCase"))
	viper.SetDefault("FUNDTRANSFERCONFIG.OauthFundTransferUrl", os.Getenv("authUrl"))
	viper.SetDefault("FUNDTRANSFERCONFIG.INQUIRYSTATUSFUNDTRANSFERURL", os.Getenv("inqUrl"))
	viper.SetDefault("FUNDTRANSFERCONFIG.CONSUMERID", os.Getenv("fundTranFerConsumerId"))
	viper.SetDefault("FUNDTRANSFERCONFIG.CONSUMERSECRET", os.Getenv("fundTranFerConsumerSecret"))
	viper.SetDefault("FUNDTRANSFERCONFIG.MerchantID", "ARRT")
	viper.SetDefault("FUNDTRANSFERCONFIG.FROMACCOUNTNO", "0481418100")
	viper.SetDefault("FUNDTRANSFERCONFIG.AUTH", os.Getenv("fundTranFerAuth"))
	viper.SetDefault("FUNDTRANSFERCONFIG.SENDERNAME", "AURORA TRADING CO.LTD.")
	viper.SetDefault("FUNDTRANSFERCONFIG.TYPEOFSENDER", "K")
	viper.SetDefault("FUNDTRANSFERCONFIG.INQUIRYSTATUSRETRY", 3)
	viper.SetDefault("FUNDTRANSFERCONFIG.OauthRetry", 5)

	viper.SetDefault("HTTP.TIMEOUT", "10s")
	viper.SetDefault("HTTP.MAXIDLECONN", 100)
	viper.SetDefault("HTTP.MAXIDLECONNPERHOST", 100)
	viper.SetDefault("HTTP.MAXCONNPERHOST", 100)
	//viper.SetDefault("HTTP.CERTFILE", "star_allgold_arrgx_com.crt")
	//viper.SetDefault("HTTP.KeyFile", "_.allgold.arrgx.com.key")

	viper.SetDefault("Producer.FinalTxnListener", os.Getenv("FinalTxnTopic"))
	viper.SetDefault("Producer.GoldListener", os.Getenv("GoldTopic"))
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
