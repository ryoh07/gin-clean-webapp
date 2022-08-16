package aws

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type AwsS3 struct {
	Config     *AWSS3Config
	Keys       AwsS3URLs
	Uploader   *s3manager.Uploader
	Downloader *s3manager.Downloader
}

type AwsS3URLs struct {
	Test string
}

func NewAwsS3() *AwsS3 {

	config := NewAWSS3Config()

	// The session the S3 Uploader will use
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Credentials: credentials.NewStaticCredentialsFromCreds(credentials.Value{
				AccessKeyID:     config.Aws.S3.AccessKeyID,
				SecretAccessKey: config.Aws.S3.SecretAccessKey,
			}),
			Region: aws.String(config.Aws.S3.Region),
		},
	}))

	return &AwsS3{
		Config: config,
		Keys: AwsS3URLs{
			Test: "test/images",
		},
		// Create an uploader with the session and default options
		Uploader:   s3manager.NewUploader(sess),
		Downloader: s3manager.NewDownloader(sess),
	}
}

func (a *AwsS3) S3Uploader(wb *bytes.Buffer, fileName string, extension string) (string, error) {

	if fileName == "" {
		return "", errors.New("fileName is required")
	}

	var contentType string

	switch extension {
	case "jpg":
		contentType = "image/jpeg"
	case "jpeg":
		contentType = "image/jpeg"
	case "gif":
		contentType = "image/gif"
	case "png":
		contentType = "image/png"
	default:
		return "", errors.New("this extension is invalid")
	}

	result, err := a.Uploader.Upload(&s3manager.UploadInput{
		Body:        wb,
		Bucket:      aws.String(a.Config.Aws.S3.Bucket),
		ContentType: aws.String(contentType),
		Key:         aws.String(a.Keys.Test + "/" + fileName + "." + extension),
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload file, %v", err)
	}
	u, _ := url.Parse(result.Location)
	fmt.Println(u.Path)
	return u.Path, nil
}

func (a *AwsS3) S3Downloader(path string) (*bytes.Buffer, error) {

	buf := aws.NewWriteAtBuffer([]byte{})

	// Upload the file to S3.
	_, err := a.Downloader.Download(buf, &s3.GetObjectInput{
		Bucket: aws.String(a.Config.Aws.S3.Bucket),
		Key:    aws.String("/test/images/hoge.png"),
	})
	if err != nil {
		panic(err)
	}

	return bytes.NewBuffer(buf.Bytes()), nil
}

func (a *AwsS3) TestS3Uploader(wb *bytes.Buffer, fileName string, extension string) (string, error) {

	return "test", nil
}

func (a *AwsS3) TestS3Downloader(path string) (*bytes.Buffer, error) {

	buf := aws.NewWriteAtBuffer([]byte{})

	return bytes.NewBuffer(buf.Bytes()), nil
}
