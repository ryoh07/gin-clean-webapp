package aws

import (
	"os"

	"github.com/ryoh07/gin-clean-webapp/common/env"
)

type AWSS3Config struct {
	Aws struct {
		S3 struct {
			Region          string
			Bucket          string
			AccessKeyID     string
			SecretAccessKey string
		}
	}
}

func NewAWSS3Config() *AWSS3Config {

	c := new(AWSS3Config)
	env.EnvLoad()
	// ex) アジアパシフィック (東京): ap-northeast-1
	c.Aws.S3.Region = os.Getenv("S3_RESION")
	c.Aws.S3.Bucket = os.Getenv("S3_BUCKET")
	c.Aws.S3.AccessKeyID = os.Getenv("S3_ACCESSKEYID")
	c.Aws.S3.SecretAccessKey = os.Getenv("S3_SECRETACCESSKEY")

	return c
}
