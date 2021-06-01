package bankly

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3Manager struct {
	s3       *session.Session
}

func NewS3Manager(s3 *session.Session) *S3Manager {
	return &S3Manager{s3: s3}
}

// Upload ...
func (s S3Manager) Upload(fileName string, bucket string, reader io.Reader) (*string, error) {
	uploader := s3manager.NewUploader(s.s3)

	response, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
		Body:   reader,
	})
	if err != nil {
		return nil, err
	}

	return &response.Location, nil
}

// Download ...
func (s S3Manager) Download(fileName string, bucket string, writer *os.File) (*int64, error) {
	downloader := s3manager.NewDownloader(s.s3)

	bucket = "temp.documentanalysis"
	fileName = "contbank.png"

	//// newFileName := "attachment; filename='" + fileName + ".png'"
	objectInput := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
		//// ResponseContentDisposition : &newFileName,
	}
	objectSize, err := downloader.Download(writer, objectInput)
	if err != nil {
		return nil, err
	}

	return &objectSize, nil
}

// ListBuckets ...
func (s S3Manager) ListBuckets() ([]string, error) {
	srv := s3.New(s.s3)

	result, err := srv.ListBuckets(nil)
	if err != nil {
		logrus.
			WithError(err).
			Error("unable to list buckets")
		return nil, err
	}

	list := []string{}
	for _, elem := range result.Buckets {
		logrus.
			Infof("* %s created on %s\n", aws.StringValue(elem.Name), aws.TimeValue(elem.CreationDate))
		list = append(list, *elem.Name)
	}

	return list, nil
}
