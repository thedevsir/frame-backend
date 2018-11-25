package storage

import (
	"io"
	"path"
	"strings"

	"github.com/minio/minio-go"
	"github.com/thedevsir/frame-backend/config"
	"github.com/thedevsir/frame-backend/services/errors"
)

var (
	cli      *minio.Client
	endpoint string
)

func Composer() {

	endpoint = config.MinioEndpoint

	var err error
	cli, err = minio.New(endpoint, config.MinioAccessKeyID, config.MinioSecretAccessKey, false)
	if err != nil {
		panic(err)
	}

	// Create public bucket
	bucketName := strings.Split(config.MinioBuckets, ",")
	for i := range bucketName {
		err = createBucket(bucketName[i], generateDownloadPolicy(bucketName[i]))
		if err != nil {
			panic(err)
		}
	}
}

func createBucket(bucketName, policy string) error {

	exists, err := cli.BucketExists(bucketName)
	if err != nil {
		return err
	}

	if !exists {
		err = cli.MakeBucket(bucketName, "")
		if err != nil {
			return err
		}
	}

	if policy != "" {
		err = cli.SetBucketPolicy(bucketName, policy)
		if err != nil {
			return err
		}
	}

	return nil
}

func generateDownloadPolicy(bucketName string) string {

	tmpl := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::{{bucket_name}}/*"]}]}`
	return strings.Replace(tmpl, "{{bucket_name}}", bucketName, 1)
}

func Put(name, bucketName string, r io.Reader, size int64, contentType string) error {

	_, err := cli.PutObject(bucketName, name, r, size, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return errors.ErrInternal
	}

	return nil
}

func Delete(name, bucketName string) error {

	_, err := cli.StatObject(bucketName, name, minio.StatObjectOptions{})
	if err != nil {
		minioError := minio.ToErrorResponse(err)
		if minioError.Code == "NoSuchKey" {
			return errors.ErrObjectNotFound
		}
	}

	err = cli.RemoveObject(bucketName, name)
	if err != nil {
		return errors.ErrInternal
	}

	return nil
}

func GetURL(name, bucketName string) string {

	return "http://" + path.Join(endpoint, bucketName, name)
}
