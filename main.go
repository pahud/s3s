package main

import (
	"fmt"
	//"log"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
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

type Resp map[string]interface{}

func BitlyURLShorten(urlStr string) string {
	api_key := os.Getenv("BITLY_TOKEN")
	b := bitly.New(api_key)
	shortURL, err := b.Links.Shorten(urlStr)
	if err != nil {
		log.Errorf("Failed to sign request", err)
	} else {
		log.Debug("The bitly URL is", shortURL)
	}
	return shortURL.URL
}

func SinaURLShorten(urlStr string) string {
	// http://api.weibo.com/2/short_url/shorten.json?source=2849184197&url_long=http://www.cnblogs.com
	apiurl := "http://api.weibo.com/2/short_url/shorten.json?source=2849184197"
	resp, err := http.Get(fmt.Sprintf("%s&url_long=%s", apiurl, url.QueryEscape(urlStr)))
	if err != nil {
		log.Warn(err)
		// handle error
	}
	defer resp.Body.Close()
	var body Resp
	json.NewDecoder(resp.Body).Decode(&body)
	// 	log.Info(body["urls"].([]interface{}))
	urls := body["urls"].([]interface{})
	url_short := urls[0].(map[string]interface{})["url_short"]
	//log.Info(url_short)
	return url_short.(string)
}

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
	// log.Println("The URL is", urlStr)

	// bitlyToken := os.Getenv("BITLY_TOKEN")
	// b := bitly.New(bitlyToken)
	// shortURL, err := b.Links.Shorten(urlStr)

	fmt.Println("Original URL:", urlStr)
	if !strings.HasPrefix(aws_region, "cn-") {
		shortURL := BitlyURLShorten(urlStr)
		if err != nil {
			log.Println("Failed to sign request", err)
		} else {
			fmt.Println("bitly URL:", shortURL)
		}
	} else {
		tcnURL := SinaURLShorten(urlStr)
		fmt.Println("t.cn URL:", tcnURL)
	}

}
