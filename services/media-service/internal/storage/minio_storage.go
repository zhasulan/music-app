package storage

import (
	"context"
	"net/url"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioStorage struct {
	client     *minio.Client
	bucket     string
	publicHost *url.URL
	expiry     time.Duration
}

func (s *MinioStorage) PublicURL(object string) string {
	u := s.PublicBase().ResolveReference(&url.URL{Path: "/" + s.bucket + "/" + object})
	return u.String()
}

// PublicBase returns the configured public host or falls back to http://localhost:9000.
func (s *MinioStorage) PublicBase() *url.URL {
	if s.publicHost != nil {
		return s.publicHost
	}
	u, _ := url.Parse("http://localhost:9000")
	return u
}

func NewMinioStorage(endpoint, accessKey, secretKey, bucket string, useSSL bool, publicURL string, expiry time.Duration) (*MinioStorage, error) {
	cli, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}
	var pub *url.URL
	if publicURL != "" {
		pub, err = url.Parse(publicURL)
		if err != nil {
			return nil, err
		}
	}
	return &MinioStorage{client: cli, bucket: bucket, publicHost: pub, expiry: expiry}, nil
}

func (s *MinioStorage) EnsureBucket(ctx context.Context) error {
	exists, err := s.client.BucketExists(ctx, s.bucket)
	if err != nil {
		return err
	}
	if !exists {
		return s.client.MakeBucket(ctx, s.bucket, minio.MakeBucketOptions{})
	}
	return nil
}

func (s *MinioStorage) PresignedURL(ctx context.Context, object string) (string, error) {
	reqParams := make(url.Values)
	urlStr, err := s.client.PresignedGetObject(ctx, s.bucket, object, s.expiry, reqParams)
	if err != nil {
		return "", err
	}
	if s.publicHost != nil {
		urlStr = rewriteHost(urlStr, s.publicHost)
	}
	return urlStr.String(), nil
}

func rewriteHost(u *url.URL, public *url.URL) *url.URL {
	u = u.ResolveReference(&url.URL{Path: u.Path, RawQuery: u.RawQuery})
	u.Scheme = public.Scheme
	u.Host = public.Host
	return u
}

func (s *MinioStorage) ObjectExists(ctx context.Context, object string) (bool, error) {
	_, err := s.client.StatObject(ctx, s.bucket, object, minio.StatObjectOptions{})
	if err != nil {
		if isNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func isNotFound(err error) bool {
	return err != nil && strings.Contains(err.Error(), "Not Found")
}
