package storages

import (
	"context"
	"io"
	"time"

	exceptions "notezy-backend/app/exceptions"
)

type PutOptions struct {
	ContentType string
	Metadata    map[string]string
}

type GetOptions struct {
}

type PresignOptions struct {
	Expires     time.Duration
	ContentType string
}

type Object struct {
	Key            string
	Data           []byte
	Size           int64
	ContentType    string
	ParseMediaType string
	LastModified   time.Time
	ETag           string
}

type StorageInterface interface {
	ListAllInTerminal()
	GetKey(ownerIndicator string, objectIndicator string) string
	NewObject(key string, reader io.Reader, size int64) (*Object, *exceptions.Exception)
	PutObjectByKey(ctx context.Context, key string, object *Object) *exceptions.Exception
	GetObjectByKey(ctx context.Context, key string, option *GetOptions) (io.ReadCloser, *Object, *exceptions.Exception)
	DeleteObjectByKey(ctx context.Context, key string) *exceptions.Exception
	PresignPutObjectByKey(ctx context.Context, key string, option *PresignOptions) (string, *exceptions.Exception)
	PresignGetObjectByKey(ctx context.Context, key string, option *PresignOptions) (string, *exceptions.Exception)
}
