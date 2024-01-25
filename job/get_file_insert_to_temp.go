package job

import (
	"bufio"
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"gitlab.com/prior-solution/aurora/standard-platform/common/reconcile_daily_batch/config"
	"gitlab.com/prior-solution/aurora/standard-platform/common/reconcile_daily_batch/internal/sftp"
	"go.uber.org/zap"
	"strings"
)

func GetFileInsertToTblTemp(
	cfg config.Config,
	ctx context.Context,
	logger *zap.Logger,
	s3Client *s3.S3,
	sftp *sftp.Client,
) error {

	if cfg.EnableS3 {
		file, err := s3Client.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(cfg.S3Config.BucketName),
			Key:    aws.String(cfg.S3Config.Key),
		})
		if err != nil {
			panic(err)
		}
		defer file.Body.Close()
		scanner := bufio.NewScanner(file.Body)
		for scanner.Scan() {
			line := scanner.Text()
			rows := strings.Split(line, "|")
			if len(rows) >= 9 {

				//TODO
			}

		}
	} else {

	}

	return nil
}
