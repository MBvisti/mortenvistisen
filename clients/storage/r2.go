package objectstorage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
	"go.uber.org/fx"

	appconfig "mortenvistisen/config"
)

var ErrObjectExists = errors.New("object already exists")

type CoverUpload struct {
	Filename    string
	ContentType string
	Size        int64
	Body        io.Reader
}

type CoverObject struct {
	Key string
	URL string
}

type CoverStore interface {
	CoverURL(filename string) (string, error)
	CoverKey(rawURL string) (string, bool)
	UploadCover(context.Context, CoverUpload) (CoverObject, error)
	Delete(context.Context, string) error
}

type r2API interface {
	PutObject(
		context.Context,
		*s3.PutObjectInput,
		...func(*s3.Options),
	) (*s3.PutObjectOutput, error)
	DeleteObject(
		context.Context,
		*s3.DeleteObjectInput,
		...func(*s3.Options),
	) (*s3.DeleteObjectOutput, error)
}

type R2 struct {
	client        r2API
	bucket        string
	coverPrefix   string
	publicBaseURL *url.URL
}

func NewR2(cfg appconfig.Config) (*R2, error) {
	awsCfg, err := awsconfig.LoadDefaultConfig(
		context.Background(),
		awsconfig.WithRegion("auto"),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.R2.AccessKeyID,
			cfg.R2.SecretAccessKey,
			"",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("load R2 client configuration: %w", err)
	}

	publicBaseURL, err := url.Parse(cfg.R2.PublicBaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse R2 public base URL: %w", err)
	}
	if publicBaseURL.Scheme == "" || publicBaseURL.Host == "" {
		return nil, errors.New("R2 public base URL must be an absolute URL")
	}
	if strings.Trim(cfg.R2.CoverPrefix, "/") == "" {
		return nil, errors.New("R2 cover prefix must not be empty")
	}

	endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.R2.AccountID)
	client := s3.NewFromConfig(awsCfg, func(options *s3.Options) {
		options.BaseEndpoint = aws.String(endpoint)
	})

	return &R2{
		client:        client,
		bucket:        cfg.R2.Bucket,
		coverPrefix:   strings.Trim(cfg.R2.CoverPrefix, "/"),
		publicBaseURL: publicBaseURL,
	}, nil
}

func (r *R2) coverKey(filename string) string {
	return path.Join(r.coverPrefix, filename)
}

func (r *R2) publicCoverPath(filename string) string {
	return path.Join("/", r.publicBaseURL.Path, r.coverKey(filename))
}

func (r *R2) CoverURL(filename string) (string, error) {
	if strings.TrimSpace(filename) == "" {
		return "", errors.New("cover filename must not be empty")
	}

	publicURL := *r.publicBaseURL
	publicURL.Path = r.publicCoverPath(filename)
	publicURL.RawPath = ""

	return publicURL.String(), nil
}

func (r *R2) CoverKey(rawURL string) (string, bool) {
	coverURL, err := url.Parse(rawURL)
	if err != nil || coverURL.Scheme != r.publicBaseURL.Scheme ||
		!strings.EqualFold(coverURL.Host, r.publicBaseURL.Host) ||
		coverURL.RawQuery != "" || coverURL.Fragment != "" {
		return "", false
	}

	prefix := path.Join("/", r.publicBaseURL.Path, r.coverPrefix) + "/"
	if !strings.HasPrefix(coverURL.Path, prefix) {
		return "", false
	}
	filename := strings.TrimPrefix(coverURL.Path, prefix)
	if filename == "" || strings.Contains(filename, "/") || path.Base(filename) != filename {
		return "", false
	}

	return r.coverKey(filename), true
}

func (r *R2) UploadCover(
	ctx context.Context,
	upload CoverUpload,
) (CoverObject, error) {
	key := r.coverKey(upload.Filename)
	publicURL, err := r.CoverURL(upload.Filename)
	if err != nil {
		return CoverObject{}, err
	}

	_, err = r.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:             aws.String(r.bucket),
		Key:                aws.String(key),
		Body:               upload.Body,
		ContentLength:      aws.Int64(upload.Size),
		ContentType:        aws.String(upload.ContentType),
		CacheControl:       aws.String("public, max-age=31536000, immutable"),
		ContentDisposition: aws.String("inline"),
		IfNoneMatch:        aws.String("*"),
	})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "PreconditionFailed" {
			return CoverObject{}, errors.Join(ErrObjectExists, err)
		}
		return CoverObject{}, fmt.Errorf("upload R2 object %q: %w", key, err)
	}

	return CoverObject{Key: key, URL: publicURL}, nil
}

func (r *R2) Delete(ctx context.Context, key string) error {
	_, err := r.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("delete R2 object %q: %w", key, err)
	}

	return nil
}

var Module = fx.Module(
	"object-storage-clients",
	fx.Provide(fx.Annotate(NewR2, fx.As(new(CoverStore)))),
)
