package storages

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	exceptions "notezy-backend/app/exceptions"
	logs "notezy-backend/app/logs"
	util "notezy-backend/app/util"
	constants "notezy-backend/shared/constants"
)

/* ============================== Interface & Constructor ============================== */

type InMemoryObject struct {
	Data           []byte
	ContentType    string
	ParseMediaType string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	ETag           string
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

/* ============================== Helper Functions ============================== */

func (s *inMemoryStorage) ListAllInTerminal() {
	logs.Info(s.data)
}

/* ============================== Implementations ============================== */

func (s *inMemoryStorage) GetKey(ownerIndicator string, objectIndicator string) string {
	salt := util.GetEnv("STORAGE_KEY_SALT", "")
	origin := "In-Memory-Key<" + ownerIndicator + "|" + objectIndicator + "|" + salt + ">"
	hash := sha256.Sum256([]byte(origin))
	return hex.EncodeToString(hash[:])
}

func (s *inMemoryStorage) GenerateETag(data []byte) string {
	return "In-Memory-ETag<" + string(int32(len(data))) + ">" + time.Now().String()
}

func (s *inMemoryStorage) NewObject(key string, reader io.Reader, size int64) (*Object, *exceptions.Exception) {
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

	contentTypes := strings.Split(http.DetectContentType(b), "; ")
	if len(contentTypes) == 0 {
		return nil, exceptions.Storage.FailedToDetectContentType()
	}

	eTag := s.GenerateETag(b)
	now := time.Now()

	object := Object{
		Key:            key,
		Data:           b,
		Size:           actualSize,
		ContentType:    contentTypes[0],
		ParseMediaType: "charset=utf-8",
		LastModified:   now,
		ETag:           eTag,
	}
	if len(contentTypes) > 1 {
		object.ParseMediaType = contentTypes[1]
	}

	return &object, nil
}

func (s *inMemoryStorage) PutObjectByKey(ctx context.Context, key string, object *Object) *exceptions.Exception {
	s.storageMutex.Lock()
	s.data[key] = &InMemoryObject{
		Data:        object.Data,
		ContentType: object.ContentType,
		CreatedAt:   object.LastModified,
		UpdatedAt:   object.LastModified,
		ETag:        object.ETag,
	}
	s.storageMutex.Unlock()

	return nil
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

// For Testing：return fake URL
func (s *inMemoryStorage) PresignPutObjectByKey(ctx context.Context, key string, option *PresignOptions) (string, *exceptions.Exception) {
	return "storage/mock://put/" + key, nil
}

// For Testing：return localhost URL, give the frontend ability to visit
func (s *inMemoryStorage) PresignGetObjectByKey(ctx context.Context, key string, option *PresignOptions) (string, *exceptions.Exception) {
	return "localhost:" + util.GetEnv("GIN_PORT", "7777") + "/" + constants.DevelopmentBaseURL + "/" + "storage/mock/files/" + key, nil
}
