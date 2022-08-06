package files

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {
	var awsRegion, bucketName, filePath string
	flag.StringVar(&awsRegion, "r", "", "AWS region")
	flag.StringVar(&bucketName, "b", "", "AWS S3 bucket to upload to")
	flag.StringVar(&filePath, "f", "", "Path to the file to upload")
	flag.Parse()

	if awsRegion == "" || bucketName == "" || filePath == "" {
		log.Println("Required arguments have not been provided.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Create AWS session.
	session, err := session.NewSession(&aws.Config{Region: aws.String(awsRegion)})
	if err != nil {
		log.Fatalf("could not initialize new aws session: %v", err)
	}

	s3Client := s3.New(session)

	// Call the upload to s3 file
	err = uploadFileToS3(s3Client, bucketName, filePath)
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
