package secret

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"gitlab.com/prior-solution/aurora/standard-platform/common/reconcile_daily_batch/config"
	"os"

	"github.com/shopspring/decimal"

	"strings"
)

type CommonConfig struct {
	PrivateKeyDEV  string `json:"PRIVATE_KEY_DEV"`
	PrivateKeyUAT  string `json:"PRIVATE_KEY_UAT"`
	PrivateKeyProd string `json:"PRIVATE_KEY_PROD"`
	AccessKey      string `json:"ACCESS_KEY"`
	SecretKey      string `json:"SECRET_KEY"`
	KafkaUsername  string `json:"KAFKA_USERNAME"`
	KafkaPassword  string `json:"KAFKA_PASSWORD"`
	RedisPassword  string `json:"REDIS_PASSWORD"`
	KeyFile        string `json:"KEY_FILE"`
	CertFile       string `json:"CERT_FILE"`
}

func ConfigCommonSecret(cfg *config.Config) error {
	secretName := cfg.AWSConfig.CommonSecret
	region := cfg.AWSConfig.Region
	sess := session.Must(session.NewSession())

	svc := secretsmanager.New(sess, aws.NewConfig().WithRegion(region))

	result, err := svc.GetSecretValue(&secretsmanager.GetSecretValueInput{SecretId: &secretName})
	if err != nil {
		return err
	}
	commonConfig := CommonConfig{}

	err = json.Unmarshal([]byte(*result.SecretString), &commonConfig)
	if err != nil {
		return err
	}

	cfg.Kafka.Internal.Username = commonConfig.KafkaUsername
	cfg.Kafka.Internal.Password = commonConfig.KafkaPassword
	switch strings.ToUpper(cfg.Env) {
	case "DEV":
		cfg.Secret.Private = commonConfig.PrivateKeyDEV
	case "UAT":
		cfg.Secret.Private = commonConfig.PrivateKeyUAT
	case "PROD":
		cfg.Secret.Private = commonConfig.PrivateKeyProd
	}
	if commonConfig.RedisPassword != "" {
		//TODO WITH Redis Password
		cfg.RedisConfig.Password = commonConfig.RedisPassword
	}

	certDecode, err := base64.StdEncoding.DecodeString(commonConfig.CertFile)
	if err != nil {
		return err
	}
	var keyDecode []byte
	if commonConfig.KeyFile != "" {
		keyDecode, err = base64.StdEncoding.DecodeString(commonConfig.KeyFile)
		if err != nil {
			fmt.Print("keyFile decode", err.Error())
		}
	}

		cfg.HTTP.KeyFile = keyDecode
		cfg.HTTP.CertFile = certDecode
	}

	decodedPrivateKey, err := base64.StdEncoding.DecodeString(cfg.Secret.Private)
	if err != nil {
		return err
	}
	cfg.Secret.Private = string(decodedPrivateKey)

	return nil
}

type RDSConfig struct {
	Username             string          `json:"username"`
	Password             string          `json:"password"`
	Engine               string          `json:"engine"`
	Host                 string          `json:"host"`
	Port                 decimal.Decimal `json:"port"`
	Dbname               string          `json:"dbname"`
	DbInstanceIdentifier string          `json:"dbInstanceIdentifier"`
}

func ConfigRDSSecret(cfg *config.Config) error {
	sess := session.Must(session.NewSession())
	secretName := cfg.AWSConfig.RDSSecret
	region := cfg.AWSConfig.Region

	svc := secretsmanager.New(sess, aws.NewConfig().WithRegion(region))
	// Create Secrets Manager client

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	result, err := svc.GetSecretValue(input)
	if err != nil {
		return err
	}

	// Decrypts secret using the associated KMS key.
	var secretString = *result.SecretString
	rdsConfig := RDSConfig{}
	err = json.Unmarshal([]byte(secretString), &rdsConfig)
	if err != nil {
		return err
	}
	cfg.DBConfig.Host = rdsConfig.Host
	cfg.DBConfig.Port = rdsConfig.Port.String()
	cfg.DBConfig.Username = rdsConfig.Username
	cfg.DBConfig.Password = rdsConfig.Password
	if strings.ToUpper(cfg.Env) == "DEV" {
		cfg.DBConfig.Name = rdsConfig.Dbname
	}

	return nil
	// Your code goes here.
}
