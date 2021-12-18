package utils

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"net/http"
	"os"
	"strings"
)

const (
	maxPartSize = int64(1024 * 1024 * 1024 * 3) // 3GB MAX FILE
	maxRetries  = 3
)

type WasabiS3Client struct {
	s3Client       *s3.S3
	s3Session      *session.Session
	s3BucketRegion string
}

func (awsClient *WasabiS3Client) DownloadObject(filekey string, filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		fmt.Println(err)
	}
	defer func() {
		closeFileErr := file.Close()
		if closeFileErr != nil {
			fmt.Println(closeFileErr)
		}
	}()

	downloader := s3manager.NewDownloader(awsClient.s3Session)
	downloadedBytes, downloadFileErr := downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(GetEnvVariable("AWS_BUCKET_NAME")),
		Key:    aws.String(filekey),
	})
	if downloadFileErr != nil {
		fmt.Println(downloadFileErr)
	}
	fmt.Println(downloadedBytes)
	return nil
}
func (awsClient *WasabiS3Client) UploadObject(bucketName, clinicName, filepath string, userID string, tipo string) (string, int64, error) {
	file, err := os.Open("./" + filepath)
	if err != nil {
		fmt.Printf("err opening file: %s", err)
		return "", 0, err
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	size := fileInfo.Size()
	buffer := make([]byte, size)
	fileType := http.DetectContentType(buffer)
	file.Read(buffer)

	fmt.Println("file.Name()")
	fmt.Println(file.Name())

	path := "/" + clinicName + "/" + userID + "/" + tipo + "/" + strings.Split(file.Name(), "/")[5]
	input := &s3.CreateMultipartUploadInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(path),
		ACL:         aws.String("public-read"),
		ContentType: aws.String(fileType),
	}
	resp, err := awsClient.s3Client.CreateMultipartUpload(input)
	if err != nil {
		fmt.Println(err.Error())
		return "No se pudo conectar", 0, err
	}

	var curr, partLength int64
	var remaining = size
	var completedParts []*s3.CompletedPart
	partNumber := 1
	for curr = 0; remaining != 0; curr += partLength {
		if remaining < maxPartSize {
			partLength = remaining
		} else {
			partLength = maxPartSize
		}
		completedPart, err := UploadPart(awsClient.s3Client, resp, buffer[curr:curr+partLength], partNumber)
		if err != nil {
			fmt.Println(err.Error())
			err := AbortMultipartUpload(awsClient.s3Client, resp)
			if err != nil {
				fmt.Println(err.Error())
			}
			return "No se pudo subir", 0, err
		}
		remaining -= partLength
		partNumber++
		completedParts = append(completedParts, completedPart)
	}
	completeResponse, err := CompleteMultipartUpload(awsClient.s3Client, resp, completedParts)
	if err != nil {
		fmt.Println(err.Error() + "abc")
		return "No se pudo subir", 0, err
	}
	var url *string

	url = completeResponse.Location

	return *url, size, nil
}

func (awsClient *WasabiS3Client) CreateBucket(bucketName string) error {
	result, err := awsClient.s3Client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
		ACL:    aws.String("public-read"),
	})
	fmt.Println(result)
	fmt.Println(err)
	return err
}

func (awsClient *WasabiS3Client) DeleteObject(bucket string, key *string) error {
	request := &s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    key,
	}
	fmt.Println(request)

	r, err := awsClient.s3Client.DeleteObject(request)
	if err != nil {
		return err
	}
	fmt.Println(r)
	return nil
}

func UploadPart(s3Client *s3.S3, resp *s3.CreateMultipartUploadOutput, fileBytes []byte, partNumber int) (*s3.CompletedPart, error) {
	tryNum := 1
	partInput := &s3.UploadPartInput{
		Body:          bytes.NewReader(fileBytes),
		Bucket:        resp.Bucket,
		Key:           resp.Key,
		PartNumber:    aws.Int64(int64(partNumber)),
		UploadId:      resp.UploadId,
		ContentLength: aws.Int64(int64(len(fileBytes))),
	}

	for tryNum <= maxRetries {
		uploadResult, err := s3Client.UploadPart(partInput)
		if err != nil {
			if tryNum == maxRetries {
				if aerr, ok := err.(awserr.Error); ok {
					return nil, aerr
				}
				return nil, err
			}
			tryNum++
		} else {
			return &s3.CompletedPart{
				ETag:       uploadResult.ETag,
				PartNumber: aws.Int64(int64(partNumber)),
			}, nil
		}
	}
	return nil, nil
}

func AbortMultipartUpload(s3Client *s3.S3, resp *s3.CreateMultipartUploadOutput) error {
	abortInput := &s3.AbortMultipartUploadInput{
		Bucket:   resp.Bucket,
		Key:      resp.Key,
		UploadId: resp.UploadId,
	}
	_, err := s3Client.AbortMultipartUpload(abortInput)
	if err != nil {
		return err
	}
	return nil
}

func CompleteMultipartUpload(s3Client *s3.S3, resp *s3.CreateMultipartUploadOutput, completedParts []*s3.CompletedPart) (*s3.CompleteMultipartUploadOutput, error) {
	completeInput := &s3.CompleteMultipartUploadInput{
		Bucket:   resp.Bucket,
		Key:      resp.Key,
		UploadId: resp.UploadId,
		MultipartUpload: &s3.CompletedMultipartUpload{
			Parts: completedParts,
		},
	}
	return s3Client.CompleteMultipartUpload(completeInput)
}

func New(accessKey, secretKey, endpoint, bucketRegion string) *WasabiS3Client {
	awsManager := &WasabiS3Client{}

	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Endpoint:         aws.String(endpoint),
		Region:           aws.String(bucketRegion),
		S3ForcePathStyle: aws.Bool(true),
	}
	s3Session, err := session.NewSession(s3Config)
	if err != nil {
		panic(err)
	}



	awsManager.s3Client = s3.New(s3Session, s3Config)
	awsManager.s3Session = s3Session
	awsManager.s3BucketRegion = bucketRegion

	return awsManager
}
