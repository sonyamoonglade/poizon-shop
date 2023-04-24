package uploader

import (
	"context"
	"fmt"

	"github.com/sonyamoonglade/s3-yandex-go/s3yandex"
)

type Config struct {
	Owner  string
	Bucket string
}

type Uploader struct {
	client  *s3yandex.YandexS3Client
	baseURL string
}

func NewUploader(client *s3yandex.YandexS3Client, baseURL string) *Uploader {
	return &Uploader{
		client:  client,
		baseURL: baseURL,
	}
}

type PutFileDTO struct {
	Filename,
	Destination,
	ContentType string
	Bytes []byte
}

func (f *Uploader) Put(ctx context.Context, dto PutFileDTO) error {
	return f.client.PutFileWithBytes(ctx, &s3yandex.PutFileWithBytesInput{
		ContentType: dto.ContentType,
		FileName:    dto.Filename,
		Destination: dto.Destination,
		FileBytes:   &dto.Bytes,
	})
}

func (f *Uploader) UrlToResource(filename string) string {
	return fmt.Sprintf("%s/%s", f.baseURL, filename)
}
