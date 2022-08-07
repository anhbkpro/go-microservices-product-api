package files

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {
	var awsRegion, bucketName, filePath string
	var timeout time.Duration

	flag.StringVar(&awsRegion, "r", "", "AWS region")
	flag.StringVar(&bucketName, "b", "", "AWS S3 bucket to upload to")
	flag.StringVar(&filePath, "f", "", "Path to the file to upload")
	flag.DurationVar(&timeout, "d", 0, "Upload timeout.")
	flag.Parse()

	if awsRegion == "" || bucketName == "" || filePath == "" {
		log.Println("Required arguments have not been provided.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Create AWS session. A Session should be shared where possible to take advantage of
	// configuration and credential caching.
	sess := session.Must(session.NewSession())

	svc := s3.New(sess)

	// Create a context with a timeout that will abort the upload if it takes
	// more than the passed in timeout.
	ctx := context.Background()
	var cancelFn func()
	if timeout > 0 {
		ctx, cancelFn = context.WithTimeout(ctx, timeout)
	}

	// Ensure the context is canceled to prevent leaking.
	if cancelFn != nil {
		defer cancelFn()
	}

	// Call the upload to s3 file
	err := uploadFileToS3(svc, bucketName, filePath)
	if err != nil {
		log.Fatalf("could not upload file: %v", err)
	}
}

func uploadFileToS3(s3Client *s3.S3, bucketName, filePath string) error {
	// Get fileName from path
	fileName := filepath.Base(filePath)

	// Open the file from the file path
	upFile, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("could not open local filepath [%v]: %+v", filePath, err)
	}
	defer upFile.Close()

	// Get the file info
	upFileInfo, _ := upFile.Stat()
	var fileSize int64 = upFileInfo.Size()
	fileBuffer := make([]byte, fileSize)
	upFile.Read(fileBuffer)

	// Put the file object to S3 with the file name
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(bucketName),
		Key:                  aws.String(fileName),
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader(fileBuffer),
		ContentLength:        aws.Int64(fileSize),
		ContentType:          aws.String(http.DetectContentType(fileBuffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	return err
}

func uploadFileToS3WithContext(svc *s3.S3, ctx context.Context, bucket, key string) {
	_, err := svc.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   os.Stdin,
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == request.CanceledErrorCode {
			fmt.Fprintf(os.Stdout, "upload canceled due to timeout , %v\n", err)
		} else {
			fmt.Fprintf(os.Stdout, "failed to upload object , %v\n", err)
		}
		os.Exit(1)
	}

	fmt.Printf("successfully uploaded file to %s/%s\n", bucket, key)
}
