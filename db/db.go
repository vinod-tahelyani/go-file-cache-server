package db

import (
	"fmt"
	"net/url"
	"sync"
	"time"

	customError "go-file-cache-server.example.com/error"
	"go-file-cache-server.example.com/models"
)

var cacheFileDB models.CacheFileDB
var dbLock *sync.Mutex

func InitializeDB() {
	cacheFileDB = make(models.CacheFileDB)
	dbLock = &sync.Mutex{}
}

func NewCacheFile(cacheFileRequestBody models.CacheFileRequestBody) (models.CacheFile, error) {
	_, err := url.Parse(cacheFileRequestBody.FileURL)
	if err != nil {
		return models.CacheFile{}, fmt.Errorf(customError.WrapError("error in parsing url", err))
	}
	if cacheFileRequestBody.AuthType == string(models.BASIC) || (len(cacheFileRequestBody.Username) > 0 && len(cacheFileRequestBody.Password) > 0) {
		return models.CacheFile{
			ID:           cacheFileRequestBody.ID,
			FileURL:      cacheFileRequestBody.FileURL,
			AuthType:     models.BASIC,
			Username:     cacheFileRequestBody.Username,
			Password:     cacheFileRequestBody.Password,
			LastModified: time.Now(),
			Mutex:        &sync.Mutex{},
			Status:       models.INITIALISED,
		}, nil
	}
	return models.CacheFile{
		ID:           cacheFileRequestBody.ID,
		FileURL:      cacheFileRequestBody.FileURL,
		AuthType:     models.NO_AUTH,
		Status:       models.INITIALISED,
		LastModified: time.Now(),
		Mutex:        &sync.Mutex{},
	}, nil
}

func GetCacheFileByURL(fileURL string) (models.CacheFile, bool)  {
	for _, cacheFile := range cacheFileDB {
		if cacheFile.FileURL == fileURL {
			return cacheFile, true
		}
	}
	return models.CacheFile{}, false
}

func GetCacheFileByID(ID string) (models.CacheFile, bool)  {
	cacheFile, ok := cacheFileDB[ID]
	return cacheFile, ok
}

func SetCacheFileStatus(ID string, status models.FileStatus) error {
	cacheFile, ok := cacheFileDB[ID]
	if !ok {
		return fmt.Errorf("file with ID %s not found", ID)
	}
	cacheFile.Mutex.Lock()
	defer cacheFile.Mutex.Unlock()
	cacheFile.Status = status
	return nil
}

func SetCacheFileName(ID string, fileName string) error {
	cacheFile, ok := cacheFileDB[ID]
	if !ok {
		return fmt.Errorf("file with ID %s not found", ID)
	}
	cacheFile.Mutex.Lock()
	defer cacheFile.Mutex.Unlock()
	cacheFile.FileName = fileName
	return nil
}

func SetCacheFilePath(ID string, filePath string) error {
	cacheFile, ok := cacheFileDB[ID]
	if !ok {
		return fmt.Errorf("file with ID %s not found", ID)
	}
	cacheFile.Mutex.Lock()
	defer cacheFile.Mutex.Unlock()
	cacheFile.LocalPath = filePath
	return nil
}

func AddCacheFileToDB(cacheFile models.CacheFile) error {
	dbLock.Lock()
	defer dbLock.Unlock()
	id := cacheFile.ID
	_, found := cacheFileDB[id]
	if found {
		return fmt.Errorf("file exist with url: %s", cacheFile.FileURL)
	}
	cacheFileDB[id] = cacheFile
	return nil
}

func UpdateCacheFileInDB(updatedCacheFile models.CacheFile) error {
	cacheFile, ok := cacheFileDB[updatedCacheFile.ID]
	if !ok {
		return fmt.Errorf("file with ID %s not found", updatedCacheFile.ID)
	}
	cacheFile.Mutex.Lock()
	defer cacheFile.Mutex.Unlock()
	cacheFileDB[updatedCacheFile.ID] = updatedCacheFile
	return nil
}

