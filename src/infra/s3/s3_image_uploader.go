package s3

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"mime/multipart"
	"os"
	"path"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	tipsDir     = "tips"
	tmpLocalDir = "tmp"
)

type S3ImageUploader struct {
	client *s3.Client
}

func NewS3ImageUploader(client *s3.Client) *S3ImageUploader {
	return &S3ImageUploader{
		client: client,
	}
}

func (u *S3ImageUploader) UploadWithFile(file multipart.File) (string, error) {
	bckt := "test-bucket"
	salt := uuid.New().String()
	filename := hash(salt + time.Now().String())
	filepath := path.Join(tipsDir, filename)

	ctx, cancelHdlr := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancelHdlr()
	_, err := u.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &bckt,
		Key:    &filepath,
		Body:   file,
	})
	if err != nil {
		return "", errors.Stack(fmt.Errorf("error uploading file using bucket %v and filename %v and filepath %v, err = %w", bckt, filename, filepath, err))
	}

	return filepath, nil
}

func (u *S3ImageUploader) Delete(path string) error {
	return nil
}

func (u *S3ImageUploader) UploadWithBase64(imageB64Data []byte) (string, error) {
	bckt := "test-bucket"
	salt := uuid.New().String()
	filename := hash(salt + time.Now().String())
	filepath := path.Join(tipsDir, filename)

	file, err := os.CreateTemp("", filename)
	if err != nil {
		return "", errors.Stack(fmt.Errorf("error creating filename %v %w", file, err))
	}
	// defer func() {
	// 	zap.S().Info("Closing file")
	// 	err = file.Close()
	// 	if err != nil {
	// 		zap.S().Errorf("Error closing file %w", err)
	// 	}
	// 	zap.S().Info("Removing file")
	// 	err = os.Remove(file.Name())
	// 	if err != nil {
	// 		zap.S().Errorf("Error removing file %w", err)
	// 	}
	// }()

	defer func() {
		zap.S().Info("Closing file")
		err = file.Close()
		if err != nil {
			zap.S().Errorf("Error closing file %w", err)
		}
	}()

	defer func() {
		zap.S().Info("Removing file")
		err = os.Remove(file.Name())
		if err != nil {
			zap.S().Errorf("Error removing file %w", err)
		}
	}()

	if _, err = file.Write(imageB64Data); err != nil {
		return "", errors.Stack(err)
	}
	if _, err := file.Seek(0, 0); err != nil { // バッファをリワインド
		return "", errors.Stack(err)
	}
	ctx, cancelHdlr := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancelHdlr()
	_, err = u.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &bckt,
		Key:    &filepath,
		Body:   file,
	})
	if err != nil {
		return "", errors.Stack(fmt.Errorf("error uploading file using bucket %v and filename %v and filepath %v, err = %w", bckt, filename, filepath, err))
	}

	return filepath, nil
}

func hash(s string) string {
	converted := sha256.Sum256([]byte(s))
	return hex.EncodeToString(converted[:])
}
