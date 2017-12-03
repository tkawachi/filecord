package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"os"
	"flag"
)

func main() {
	var bucket, key string
	flag.StringVar(&bucket, "b", "", "Bucket name.")
	flag.StringVar(&key, "k", "", "Object key name.")
	flag.Parse()

	download(bucket, key)
}
func download(bucket string, key string) {
	sess := session.Must(session.NewSession())
	svc := s3.New(sess)
	input := &s3.GetObjectInput{Bucket: aws.String(bucket), Key: aws.String(key)}
	if err := input.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid request, %v", err)
		os.Exit(1)
	}

	output, err := svc.GetObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == request.CanceledErrorCode {
			fmt.Fprintf(os.Stderr, "get canceled due to timeout, %v\n", err)
		} else {
			fmt.Fprintf(os.Stderr, "failed to get object, %v\n", err)
		}
		os.Exit(1)
	}
	defer output.Body.Close()
	fmt.Println("Object size:", aws.Int64Value(output.ContentLength))
}
