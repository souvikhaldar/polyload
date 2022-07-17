package uploader

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/spf13/viper"
)

type Aws struct {
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
}

func NewAws() *Aws {
	// The session the S3 Uploader will use
	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(viper.GetString("aws.s3_region")),
			Credentials: credentials.NewStaticCredentials(
				viper.GetString("aws.access_key_id"),
				viper.GetString("aws.secret_access_key"),
				"", // a token will be created when the session it's used.
			),
		})
	if err != nil {
		log.Fatal(err)
	}

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)
	// Create a downloader with the session and default options
	downloader := s3manager.NewDownloader(sess)
	aws := new(Aws)
	aws.uploader = uploader
	aws.downloader = downloader

	return aws
}
func (a *Aws) UploadFile(fileName string, data []byte) error {
	// Upload the file to S3.
	r := bytes.NewReader(data)

	_, err := a.uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(viper.GetString("aws.s3_bucket_name")),
		Key:    aws.String(fileName),
		Body:   r,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file, %v", err)
	}
	return nil
}
func (a *Aws) DownloadFile(fileName string) ([]byte, error) {
	// Create a file to write the S3 Object contents to.
	f, err := os.Create(viper.GetString("download_dir") + fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to create file %q, %v", fileName, err)
	}

	// Write the contents of S3 Object to the file
	n, err := a.downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(viper.GetString("aws.s3_bucket_name")),
		Key:    aws.String(fileName),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to download file, %v", err)
	}
	fileContents, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("Error in reading file: ", err)
	}
	fmt.Printf("file downloaded, %d bytes\n", n)
	return fileContents, nil
}
