package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/zpnk/go-bitly"
)

var (
	aws_profile = "default"
	aws_region  = "ap-northeast-1"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("usage: %s <bucket> <filename>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	bucket := os.Args[1]
	filename := os.Args[2]

	if profile := os.Getenv("AWS_PROFILE"); profile != "" {
		aws_profile = profile
	}
	if region := os.Getenv("AWS_DEFAULT_REGION"); region != "" {
		aws_region = region
	}

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Failed to open file", filename, err)
		os.Exit(1)
	}
	defer file.Close()

	//select Region to use.
	// conf := aws.Config{
	// 	Region: aws.String("ap-northeast-1"),
	// }
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config:  aws.Config{Region: aws.String(aws_region)},
		Profile: aws_profile,
	}))
	// sess := session.New(&conf)
	svc := s3manager.NewUploader(sess)
	s3svc := s3.New(sess)

	fmt.Println("Uploading file to S3...")
	result, err := svc.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filepath.Base(filename)),
		Body:   file,
	})
	if err != nil {
		fmt.Println("error", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully uploaded %s to %s\n", filename, result.Location)
	// presign
	req, _ := s3svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filepath.Base(filename)),
	})
	urlStr, err := req.Presign(7 * 86400 * time.Second)
	if err != nil {
		log.Println("Failed to sign request", err)
	}
	log.Println("The URL is", urlStr)
	bitlyToken := os.Getenv("BITLY_TOKEN")
	b := bitly.New(bitlyToken)
	shortURL, err := b.Links.Shorten(urlStr)
	if err != nil {
		log.Println("Failed to sign request", err)
	} else {
		fmt.Println("LongURL:", shortURL.LongURL)
		fmt.Println("ShortURL:", shortURL.URL)
	}

}
