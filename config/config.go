package config

import (
	"encoding/base64"
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
	viper.SetDefault("LOG.LEVEL", os.Getenv("logLevel"))
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
	viper.SetDefault("FUNDTRANSFERCONFIG.FROMACCOUNTNO", os.Getenv("fundTranFerFrom"))
	viper.SetDefault("FUNDTRANSFERCONFIG.AUTH", os.Getenv("fundTranFerAuth"))
	viper.SetDefault("FUNDTRANSFERCONFIG.SENDERNAME", "AURORA TRADING CO.LTD.")
	viper.SetDefault("FUNDTRANSFERCONFIG.TYPEOFSENDER", "K")
	viper.SetDefault("FUNDTRANSFERCONFIG.INQUIRYSTATUSRETRY", os.Getenv("inqRetry"))
	viper.SetDefault("FUNDTRANSFERCONFIG.OauthRetry", os.Getenv("authRetry"))

	viper.SetDefault("HTTP.TIMEOUT", "10s")
	viper.SetDefault("HTTP.MAXIDLECONN", 100)
	viper.SetDefault("HTTP.MAXIDLECONNPERHOST", 100)
	viper.SetDefault("HTTP.MAXCONNPERHOST", 100)
	//certStr := `LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUdLVENDQlJHZ0F3SUJBZ0lRQTQxREZ6Z252dnV3V3Y0N2t4UTVlekFOQmdrcWhraUc5dzBCQVFzRkFEQmUKTVFzd0NRWURWUVFHRXdKVlV6RVZNQk1HQTFVRUNoTU1SR2xuYVVObGNuUWdTVzVqTVJrd0Z3WURWUVFMRXhCMwpkM2N1WkdsbmFXTmxjblF1WTI5dE1SMHdHd1lEVlFRREV4UlVhR0YzZEdVZ1ZFeFRJRkpUUVNCRFFTQkhNVEFlCkZ3MHlNekEzTVRNd01EQXdNREJhRncweU5EQTNNVEl5TXpVNU5UbGFNQjR4SERBYUJnTlZCQU1NRXlvdVlXeHMKWjI5c1pDNWhjbkpuZUM1amIyMHdnZ0VpTUEwR0NTcUdTSWIzRFFFQkFRVUFBNElCRHdBd2dnRUtBb0lCQVFDbgpzenA1RFhLTkI2OFpOd1RDdElaVFEzbzJWRlluaUdKNjNMUTJRMXRDSW5iZFEzVkRuQVpRNFAwVC9YZ0pPSkU2CnQxNE9La1Z1ZmdrZTJnMzJ0ZlJpR2tXTWU0bGhzUjlpSjV0aWhseEpRNS84WU9CeklzK1ZSNnVwZk5zR0hrOEoKcjR2VmdOV2NCa0JEc2NYYitnT2VycVV6UENYKzkrK2tLOUplZElNQXl2SE9mYXhTV2FqNUYwS1krSzQycks4TQpGaEpFTDR5VmM3Z0syNFozZG5tU2xBVm5IU3ovVU5RS00rbGQ3RmpCeFh2TDYvRWZ2M3lHSHJtY2tmWmRQVVpECnowaGszTDIyRzJSQS93ckhuSDJWYWNWSENxKy9saEJDV0xadmFzcUVUQVNkYzJsZTVwc2piVkQ5MUF3YStOTHkKV3Q1NUhwMVl2TEoza2RWc1NGZERBZ01CQUFHamdnTWhNSUlESFRBZkJnTlZIU01FR0RBV2dCU2xqUDR5ek9zUApMTlFaeGdpNEFDU0lYY1BGdHpBZEJnTlZIUTRFRmdRVWJyMGM0ZkZrQ0M3R0pNQ2JUdFBMZlhWb2lWb3dNUVlEClZSMFJCQ293S0lJVEtpNWhiR3huYjJ4a0xtRnljbWQ0TG1OdmJZSVJZV3hzWjI5c1pDNWhjbkpuZUM1amIyMHcKRGdZRFZSMFBBUUgvQkFRREFnV2dNQjBHQTFVZEpRUVdNQlFHQ0NzR0FRVUZCd01CQmdnckJnRUZCUWNEQWpBNwpCZ05WSFI4RU5EQXlNRENnTHFBc2hpcG9kSFJ3T2k4dlkyUndMblJvWVhkMFpTNWpiMjB2VkdoaGQzUmxWRXhUClVsTkJRMEZITVM1amNtd3dQZ1lEVlIwZ0JEY3dOVEF6QmdabmdRd0JBZ0V3S1RBbkJnZ3JCZ0VGQlFjQ0FSWWIKYUhSMGNEb3ZMM2QzZHk1a2FXZHBZMlZ5ZEM1amIyMHZRMUJUTUhBR0NDc0dBUVVGQndFQkJHUXdZakFrQmdncgpCZ0VGQlFjd0FZWVlhSFIwY0RvdkwzTjBZWFIxY3k1MGFHRjNkR1V1WTI5dE1Eb0dDQ3NHQVFVRkJ6QUNoaTVvCmRIUndPaTh2WTJGalpYSjBjeTUwYUdGM2RHVXVZMjl0TDFSb1lYZDBaVlJNVTFKVFFVTkJSekV1WTNKME1Ba0cKQTFVZEV3UUNNQUF3Z2dGOUJnb3JCZ0VFQWRaNUFnUUNCSUlCYlFTQ0FXa0Jad0IxQU83TjBHVFYyeHJPeFZ5MwpuYlRORTZJeWgwWjh2T3pldzFGSVdVWnhIN1diQUFBQmlVNVFQY0lBQUFRREFFWXdSQUlnUmc5WVN6QmU0OWVqCk9QdnAzVVhkSGl0ZUx5NTVwUElQeFJRZjNnWEoyVjRDSUNRUThNVFBUMmpsRHJhMENEL0J4Tzl4dGhrQ2lnSUUKS1ZUTUxqdDVYQkJhQUhZQTJyYS9heisxdGlLZm04SzdYR3ZvY0pGeGJMdFJoSVUwdmFROU1FalgrNnNBQUFHSgpUbEErSlFBQUJBTUFSekJGQWlFQThwaDNyY1Uvb2hPSUhQME5RYTFPZXpQd3B1OWl0SGJKWU1TNllWM1JmaW9DCklDK01SZVphM1N4eS93RjY0Wnk5cTF2UmxqYmhQVVhoZnJzZmU2MEVqN3ZOQUhZQU8xTjNkVDR0dVlCT2l6QmIKQnY1QU8yZllUOFAweDcwQURTMXliK0g2MUJjQUFBR0pUbEErQ2dBQUJBTUFSekJGQWlFQSs2UnBXV082bXRpYgpQN2MyS2JjbU1LcWtsR2YxUUlzVWhCaGJLQVBTVEswQ0lFSTdWSEhWaVN4T3U3UEJRR204a2NYK25tUks5NllICnFkT3BhSTNQZUdwUk1BMEdDU3FHU0liM0RRRUJDd1VBQTRJQkFRQWZVbmllRGY3c3JuY2pxaytFSzRzeHFPR1cKaXZKMjNSTmFtVlpCb3Y1R29BSVBnODg2bjlYYm5WajlMdDRMQ0tQM3NOOTRRTHRUbmlXd0JVZnprQUhVbGZrUwpPVkl3QXBOVEdWRTJaeU0rZmJMVlFMMkVDZFN5Y1JSOVd2N1duKzdydE1kQVFmQytWRVZHdE1CdC9Tc2FOTHp5CmgwbTR6VXRtN3VpYTgxcGQrZDRUczRQR1dVZllYV0Y1MC9PVkVxaVFrUjJaN1MwQlBPbDVJV3BTbnNlcEtLZGYKeEgydEgvcW92MFRQWlQxWmhxZHB6N1hCK3ZQOGJOWHRqc1N5NFVlSUhZOHhmRjJUZmNLN0lYQkVBdWNFU200RwpFVmJEeU9KL0JUdUlCaFRac280NUxKdWxrdS9vNEIxUUpGbFR2c2ZzcTVvcnI0dmRCMHJqOTE4K1pvOVkKLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQ==`
	certStr := `LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUVKRENDQXd5Z0F3SUJBZ0lTQkd4bmVYN0Fzb3pDd0greHFlMVZSN0RlTUEwR0NTcUdTSWIzRFFFQkN3VUEKTURJeEN6QUpCZ05WQkFZVEFsVlRNUll3RkFZRFZRUUtFdzFNWlhRbmN5QkZibU55ZVhCME1Rc3dDUVlEVlFRRApFd0pTTXpBZUZ3MHlOREF6TVRJd09EUTVOVGRhRncweU5EQTJNVEF3T0RRNU5UWmFNQjB4R3pBWkJnTlZCQU1NCkVpb3VjSFZpYkdsakxtRnljbWQ0TG1OdmJUQlpNQk1HQnlxR1NNNDlBZ0VHQ0NxR1NNNDlBd0VIQTBJQUJLdW8KNWZFdU00QTV4K3N4a3B4VGNBVXJqUEtuTktwZUxnYzQzanRsamZlemJLeVE3a1RXMWZSRjl3ZG5DM1hzOGJtWgp4bi9mMk9VdnVTTDdKQjNYSWttamdnSVNNSUlDRGpBT0JnTlZIUThCQWY4RUJBTUNCNEF3SFFZRFZSMGxCQll3CkZBWUlLd1lCQlFVSEF3RUdDQ3NHQVFVRkJ3TUNNQXdHQTFVZEV3RUIvd1FDTUFBd0hRWURWUjBPQkJZRUZLbTkKVHhPU3FpYldUUHlvOWdRaDlGY21BT29oTUI4R0ExVWRJd1FZTUJhQUZCUXVzeGUzV0ZiTHJsQUpRT1lmcjUyTApGTUxHTUZVR0NDc0dBUVVGQndFQkJFa3dSekFoQmdnckJnRUZCUWN3QVlZVmFIUjBjRG92TDNJekxtOHViR1Z1ClkzSXViM0puTUNJR0NDc0dBUVVGQnpBQ2hoWm9kSFJ3T2k4dmNqTXVhUzVzWlc1amNpNXZjbWN2TUIwR0ExVWQKRVFRV01CU0NFaW91Y0hWaWJHbGpMbUZ5Y21kNExtTnZiVEFUQmdOVkhTQUVEREFLTUFnR0JtZUJEQUVDQVRDQwpBUUlHQ2lzR0FRUUIxbmtDQkFJRWdmTUVnZkFBN2dCMUFEdFRkM1UrTGJtQVRvc3dXd2IrUUR0bjJFL0Q5TWU5CkFBMHRjbS9oK3RRWEFBQUJqaklSbHkwQUFBUURBRVl3UkFJZ1dpSFBpbWUxZ3VWLzVBQzhhY1ZXaVpRU0ZtVy8KTDQzVlBDZWRwdkg2K1VVQ0lEWnFHRFU3VUw3R2lra1I3bGJYMkZvNjh3TUd4Y1J1V3N4OE1Vc3hhMnowQUhVQQpTTERqYTlxbVJ6UVA1V29DK3AwdzZ4eFNBY3RXM1N5QjJidS9xem5ZaEhNQUFBR09NaEdaS0FBQUJBTUFSakJFCkFpQWZSSjg2KzBLcHgvVkZEYVFBeGw0aGhiQ2lLUUpmcy9YVTJ0MHBKTlJlSkFJZ2RtQ05sMGtDWk1kSlhONi8KR1N4bTc2bE5teE5MK2NOT3pnbzdsdStiMzJ3d0RRWUpLb1pJaHZjTkFRRUxCUUFEZ2dFQkFDNVU2SzMzQ29CVwpmY1h3dUFXemk2ZUVuSTRlVEorbzdDUVl0RkJwQ25tTURoZHpuN1hlTGx3YlBUZzZvdTQrK0l5ekJxNjAyRWh2CjZFUFR5emlEZDZpS2dIOVRjN1J1Z2JwRE5mV0JuM3dSNmV3NEZkSENrL05pR2ZTR0t6WVI0K2FKMSsweWtZenEKNTlQbWlFQ0EyVnF3MkNuS2FvTy9qdTJmREo1d09zZFpWT2JrUkhQSVpXbVR1d2EzVHovVVozWjdrWEVlODhZOQpkVDdBeC8xOFNtKy80U0Fna3EvRVVxVXp1VkY2OWttM1Qyd05qVzdxV3g1bUtWYVoxcTE2dlVXSWdzSXlBTnJyCkRJM1VBTnVjUlFpNUkvcE0vWEs4NWQycTdVbW9PbGxVREQwdTFFdWJYb3JMMjNla281VlNOUEtsUTN3MWduUy8KT3VuOVAra1cwVFU9Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0KLS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUZGakNDQXY2Z0F3SUJBZ0lSQUpFckNFclBEQmluVS9iV0xpV25YMW93RFFZSktvWklodmNOQVFFTEJRQXcKVHpFTE1Ba0dBMVVFQmhNQ1ZWTXhLVEFuQmdOVkJBb1RJRWx1ZEdWeWJtVjBJRk5sWTNWeWFYUjVJRkpsYzJWaApjbU5vSUVkeWIzVndNUlV3RXdZRFZRUURFd3hKVTFKSElGSnZiM1FnV0RFd0hoY05NakF3T1RBME1EQXdNREF3CldoY05NalV3T1RFMU1UWXdNREF3V2pBeU1Rc3dDUVlEVlFRR0V3SlZVekVXTUJRR0ExVUVDaE1OVEdWMEozTWcKUlc1amNubHdkREVMTUFrR0ExVUVBeE1DVWpNd2dnRWlNQTBHQ1NxR1NJYjNEUUVCQVFVQUE0SUJEd0F3Z2dFSwpBb0lCQVFDN0FoVW96UGFnbE5NUEV1eU5WWkxEK0lMeG1hWjZRb2luWFNhcXRTdTV4VXl4cjQ1citYWElvOWNQClI1UVVWVFZYako2b29qa1o5WUk4UXFsT2J2VTd3eTdiamNDd1hQTlpPT2Z0ejJud1dnc2J2c0NVSkNXSCtqZHgKc3hQbkhLemhtKy9iNUR0RlVrV1dxY0ZUempUSVV1NjFydTJQM21CdzRxVlVxN1p0RHBlbFFEUnJLOU84WnV0bQpOSHo2YTR1UFZ5bVorREFYWGJweWIvdUJ4YTNTaGxnOUY4Zm5DYnZ4Sy9lRzNNSGFjVjNVUnVQTXJTWEJpTHhnClozVm1zL0VZOTZKYzVsUC9Pb2kyUjZYL0V4anFtQWwzUDUxVCtjOEI1ZldtY0JjVXIyT2svNW16azUzY1U2Y0cKL2tpRkhhRnByaVYxdXhQTVVnUDE3VkdoaTlzVkFnTUJBQUdqZ2dFSU1JSUJCREFPQmdOVkhROEJBZjhFQkFNQwpBWVl3SFFZRFZSMGxCQll3RkFZSUt3WUJCUVVIQXdJR0NDc0dBUVVGQndNQk1CSUdBMVVkRXdFQi93UUlNQVlCCkFmOENBUUF3SFFZRFZSME9CQllFRkJRdXN4ZTNXRmJMcmxBSlFPWWZyNTJMRk1MR01COEdBMVVkSXdRWU1CYUEKRkhtMFdlWjd0dVhrQVhPQUNJaklHbGoyNlp0dU1ESUdDQ3NHQVFVRkJ3RUJCQ1l3SkRBaUJnZ3JCZ0VGQlFjdwpBb1lXYUhSMGNEb3ZMM2d4TG1rdWJHVnVZM0l1YjNKbkx6QW5CZ05WSFI4RUlEQWVNQnlnR3FBWWhoWm9kSFJ3Ck9pOHZlREV1WXk1c1pXNWpjaTV2Y21jdk1DSUdBMVVkSUFRYk1Ca3dDQVlHWjRFTUFRSUJNQTBHQ3lzR0FRUUIKZ3Q4VEFRRUJNQTBHQ1NxR1NJYjNEUUVCQ3dVQUE0SUNBUUNGeWs1SFBxUDNoVVNGdk5WbmVMS1lZNjExVFI2VwpQVE5sY2xRdGdhRHF3KzM0SUw5ZnpMZHdBTGR1Ty9aZWxON2tJSittNzR1eUErZWl0Ulk4a2M2MDdUa0M1M3dsCmlrZm1aVzQvUnZUWjhNNlVLKzVVemhLOGpDZEx1TUdZTDZLdnpYR1JTZ2kzeUxnamV3UXRDUGtJVno2RDJRUXoKQ2tjaGVBbUNKOE1xeUp1NXpsenlaTWpBdm5uQVQ0NXRSQXhla3JzdTk0c1E0ZWdkUkNuYldTRHRZN2toK0JJbQpsSk5Yb0IxbEJNRUtJcTRRRFVPWG9SZ2ZmdURnaGplMVdyRzlNTCtIYmlzcS95Rk9Hd1hEOVJpWDhGNnN3Nlc0CmF2QXV2RHN6dWU1TDNzejg1SytFQzRZL3dGVkROdlpvNFRZWGFvNlowZitsUUtjMHQ4RFFZemsxT1hWdThycDIKeUpNQzZhbExiQmZPREFMWnZZSDduN2RvMUFabHM0STlkMVA0am5rRHJRb3hCM1VxUTloVmwzTEVLUTczeEYxTwp5SzVHaEREWDhvVmZHS0Y1dStkZWNJc0g0WWFUdzdtUDNHRnhKU3F2MyswbFVGSm9pNUxjNWRhMTQ5cDkwSWRzCmhDRXhyb0wxKzdtcnlJa1hQZUZNNVRnTzlyMHJ2WmFCRk92VjJ6MGdwMzVaMCtMNFdQbGJ1RWpOL2x4UEZpbisKSGxVanI4Z1JzSTNxZkpPUUZ5LzlyS0lKUjBZLzhPbXd0LzhvVFdneTFtZGVIbW1qazdqMW5Zc3ZDOUpTUTZadgpNbGRsVFRLQjN6aFRoVjErWFdZcDZyamQ1SlcxemJWV0VrTE54RTdHSlRoRVVHM3N6Z0JWR1A3cFNXVFVUc3FYCm5MUmJ3SE9vcTdoSHdnPT0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQ==`
	certDecode, _ := base64.StdEncoding.DecodeString(certStr)
	viper.SetDefault("HTTP.CERTFILE", certDecode)

	//keyStr := `LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCk1JSUV2Z0lCQURBTkJna3Foa2lHOXcwQkFRRUZBQVNDQktnd2dnU2tBZ0VBQW9JQkFRQ25zenA1RFhLTkI2OFoKTndUQ3RJWlRRM28yVkZZbmlHSjYzTFEyUTF0Q0luYmRRM1ZEbkFaUTRQMFQvWGdKT0pFNnQxNE9La1Z1ZmdrZQoyZzMydGZSaUdrV01lNGxoc1I5aUo1dGlobHhKUTUvOFlPQnpJcytWUjZ1cGZOc0dIazhKcjR2VmdOV2NCa0JECnNjWGIrZ09lcnFVelBDWCs5KytrSzlKZWRJTUF5dkhPZmF4U1dhajVGMEtZK0s0MnJLOE1GaEpFTDR5VmM3Z0sKMjRaM2RubVNsQVZuSFN6L1VOUUtNK2xkN0ZqQnhYdkw2L0VmdjN5R0hybWNrZlpkUFVaRHowaGszTDIyRzJSQQovd3JIbkgyVmFjVkhDcSsvbGhCQ1dMWnZhc3FFVEFTZGMybGU1cHNqYlZEOTFBd2ErTkx5V3Q1NUhwMVl2TEozCmtkVnNTRmREQWdNQkFBRUNnZ0VBRENnN0VmRituMml5TWVMQ0xwZEZzWjJQcTRhYnBFd0h6NTVXVmlTMTVlcDMKc1h5bGNKeEwvT3NDamNOdlEwUGRpMk1scDJNN0cxSjV1TW5YLzAyYmhNMGd3NWxsRVRiMDdubXVrd3JvZjhzdQpPdTZPOXVuTUlLZE1jNElBb3NYcHR1c0orUlZZNXZHeEVQYy9QNzQxS3ZqQU15R21JNEMzMTYveGxUVmZGZHlDCmRVQ2ZCV3FXWlVLUE9PeDZIWVVBSzB3TGwvdDNTQWJwWGxNdVNkL3k1b0liS0FMWG5rN1FQbHdFYzg2ZCs1VVYKelFDREJmN2hyTzJRUXN5RitiYzFWa0R0TnExU0pGSGNieE9lczZpUHA1cElhVnQ2OW1lRDdCcFBMVmFDU2hERwpoWURGUUtSSGJ4WFYxMkU0cGduZlR4WDg5cXprYVNkSTNNTXR6Qms2aVFLQmdRRFM1UG5FbE9JZCsxMGNOc0JTCkhUOHFFWUZlUitvcXp0VFpIam1SUVlxOWpLY3N5Y2JlN3lrQ0txTG9wR0M1QXhTOWZhMjhOL0FUWkMzL3hiRlMKSmpUWW8yRVZyREs5R0g5U0oxWSs1aWUwTUNtNi9HSDNzTzUwQ2p4RDlSUm9lVEtPMnFxRytIWjdZSXBMajltYgoyc3dDbGk1bUUrTTdINTlpT2U3WnR0Qkgxd0tCZ1FETGtULy9Wc2xMVHlXNnBRMUFQZWFPN21SRWpaZGcxVklECkEzTVpvZitGazRmSGEyTW1MbW9zcUZKbWhSejFkbTNhZTJGa0lrZGZjWjc1Z28vR1QyNmtaclVIbFFTeXFtOEoKbHdzSnE0RTNoanJycHNCcEVnWU0yeVpZdlpuUHFrUjdNQ2I2V3lYUWxrSEFqaGJqSzRKNFExRkZIZGdnTnFiOApFb2NZYWV4T2RRS0JnUUNYTUU0YTd2MDNyMGR0L1paY2g5a0xpS2NzOXZOYUl4TVdZQU8zTGJ5UDdQREFQQnRWCklURk4rMUQwNVRydUI5WnJqbGpwMFZSTUlvcVRqWjkwbkMxUWpiZ0ErSVViYVIrRnZ1dW1oZ3M2c3ppSGMzMnMKTzJ5SFJmczBZTk56bmtkdmdEVzJNeE9GbVkwclpJSUZxSktPM0NtQlJvcWxqU01QSVNjcGIxVGIyd0tCZ1FDcQorTFI2ZldhVk5NVm9iR1dqdGhtbHBEMWNnbHRJdmdHaWZFdzRsQ0hyQzR5M2hjOEJhMnhMVTVmWmVTVm9WKzVOCjJPQmtYSkg3Ykk5cjJpZHRGSnZGd21sN0U4S2RXSjNudlE5Tk1ObFhUQXJDand2OWMyRFhmVnhJbmYzSU42WksKbkpld0g4dXowKzhualc5Vm50NTJxWHRoaEg1WUYrN0p1Ym56WEV0WFhRS0JnQlNEMk42OVlESVdzbUtlZUpjVQp6Z0JhazByV01hVWVpWUphZVdTTytNNGd4TnltYWZBeDdxWFQwYVNEZkpSS1l1bkd6Uyt6TERrOGg2TWhwbEZZCjJSdmt2WS9XMUNlOFdISkpBV3JRNFNHMC9aNTJ3YldVQ3Z2ckt4anh2Skc1S2VSbWxtUjBKNHNoTkQ0OVBSQVUKWmRONGVYZGxWQS9QSmJyTWtaaE1ObHdkCi0tLS0tRU5EIFBSSVZBVEUgS0VZLS0tLS0=`
	keyStr := `LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCk1JR0hBZ0VBTUJNR0J5cUdTTTQ5QWdFR0NDcUdTTTQ5QXdFSEJHMHdhd0lCQVFRZ3FuUXdKd3V5Z3B0b3l0ZGQKT1cxNWMyQnZiNjBZRy9QYm9KSENqUjA2Ym9LaFJBTkNBQVNycU9YeExqT0FPY2ZyTVpLY1UzQUZLNHp5cHpTcQpYaTRIT040N1pZMzNzMnlza081RTF0WDBSZmNIWnd0MTdQRzVtY1ovMzlqbEw3a2kreVFkMXlKSgotLS0tLUVORCBQUklWQVRFIEtFWS0tLS0t`
	keyDecode, _ := base64.StdEncoding.DecodeString(keyStr)
	viper.SetDefault("HTTP.KeyFile", keyDecode)

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
