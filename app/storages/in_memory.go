package storages

import (
	"bytes"
	"context"
	"io"
	"sync"
	"time"

	exceptions "notezy-backend/app/exceptions"
	constants "notezy-backend/shared/constants"
)

/* ============================== Interface & Constructor ============================== */

type InMemoryObject struct {
	Data        []byte
	ContentType string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ETag        string
}

type inMemoryStorage struct {
	storageMutex sync.RWMutex
	data         map[string]*InMemoryObject
}

func newInMemoryStorage() StorageInterface {
	return &inMemoryStorage{
		data: map[string]*InMemoryObject{},
	}
}

var InMemoryStorage = newInMemoryStorage()

/* ============================== Helper methods ============================== */

func (s *inMemoryStorage) GenerateKey(ownerIndicator string, objectIndicator string) string {
	return "In-Memory-Key<" + ownerIndicator + "|" + objectIndicator + "|" + ">"
}

func (s *inMemoryStorage) GenerateETag(data []byte) string {
	return "In-Memory-ETag<" + string(int32(len(data))) + ">" + time.Now().String()
}

/* ============================== Implementations ============================== */

func (s *inMemoryStorage) PutObjectByKey(ctx context.Context, key string, reader io.Reader, size int64) (*Object, *exceptions.Exception) {
	if size > constants.MaxInMemoryStorageFileSize {
		return nil, exceptions.Storage.ObjectTooLarge(size, constants.MaxInMemoryStorageFileSize)
	}

	limitReader := io.LimitReader(reader, constants.MaxInMemoryStorageFileSize+1)
	b, err := io.ReadAll(limitReader)
	if err != nil {
		return nil, exceptions.Storage.FailedToReadObjectBytes()
	}

	actualSize := int64(len(b))
	if actualSize > constants.MaxInMemoryStorageFileSize {
		return nil, exceptions.Storage.ObjectTooLarge(actualSize, constants.MaxInMemoryStorageFileSize)
	}

	eTag := s.GenerateETag(b)
	now := time.Now()

	s.storageMutex.Lock()
	s.data[key] = &InMemoryObject{
		Data:        b,
		ContentType: "application/octet-stream",
		CreatedAt:   now,
		UpdatedAt:   now,
		ETag:        eTag,
	}
	s.storageMutex.Unlock()

	return &Object{
		Key:          key,
		Size:         actualSize,
		ContentType:  "application/octec-stream",
		LastModified: now,
		ETag:         eTag,
	}, nil
}

func (s *inMemoryStorage) GetObjectByKey(ctx context.Context, key string, option *GetOptions) (io.ReadCloser, *Object, *exceptions.Exception) {
	s.storageMutex.RLock()
	object, ok := s.data[key]
	s.storageMutex.RUnlock()
	if !ok {
		return nil, nil, exceptions.Storage.FailedToGetObject(key)
	}

	rc := io.NopCloser(bytes.NewReader(object.Data))
	metadata := &Object{
		Key:          key,
		Data:         object.Data,
		Size:         int64(len(object.Data)),
		ContentType:  object.ContentType,
		LastModified: object.UpdatedAt,
		ETag:         object.ETag,
	}

	return rc, metadata, nil
}

func (s *inMemoryStorage) DeleteObjectByKey(ctx context.Context, key string) *exceptions.Exception {
	s.storageMutex.Lock()
	defer s.storageMutex.Unlock()
	if _, ok := s.data[key]; !ok {
		return exceptions.Storage.FailedToGetObject(key)
	}
	delete(s.data, key)
	return nil
}

func (s *inMemoryStorage) PresignPutObjectByKey(ctx context.Context, key string, option *PresignOptions) (string, *exceptions.Exception) {
	// For Testing：return fake URL
	return "storage/mock://put/" + key, nil
}

func (s *inMemoryStorage) PresignGetObjectByKey(ctx context.Context, key string, option *PresignOptions) (string, *exceptions.Exception) {
	// For Testing：return localhost URL, give the frontend ability to visit
	return "http://localhost:8080/storage/mock/files/" + key, nil
}
