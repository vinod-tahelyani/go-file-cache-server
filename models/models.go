package models

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	customError "go-file-cache-server.example.com/error"
)

type FileStatus int

const (
	INITIALISED FileStatus = iota
	DOWNLOADING
	ERROR_DOWLOADING
	DOWNLOADED
)

type AuthType string

const (
	NO_AUTH AuthType = ""
	BASIC   AuthType = "Basic"
	BEARER  AuthType = "Bearer"
)

type CacheFile struct {
	ID string `json:id`
	FileName     string     `json:"fileName"`
	FileURL      string     `json:"fileURL"`
	AuthType     AuthType   `json:authType`
	Username     string     `json:"username"`
	Password     string     `json:"password"`
	LocalPath    string     `json:"localPath"`
	Status       FileStatus `json:"status"`
	LastModified time.Time  `json:"lastModified"`
	Mutex *sync.Mutex
}

type CacheFileRequestBody struct {
	ID string `json:id`
	FileURL  string `json:"fileURL"`
	AuthType string `json:authType`
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

type CacheFileResponseBody struct {
	ID string `json:id`
	Status FileStatus `json:"status"`
	Message string `json:"message"`
}

type CacheFileDB map[string]CacheFile

func InitialiseNewCacheFileDB(seedCacheFilePath string) (CacheFileDB, error) {
	if seedCacheFilePath == "" {
		return make(map[string]CacheFile), nil
	}
	file, err := os.Open(seedCacheFilePath)
	if err != nil {
		return map[string]CacheFile{}, fmt.Errorf(customError.WrapError("error in opening file " + seedCacheFilePath, err))
	}
	bs, err := io.ReadAll(file)
	var cacheFileDB map[string]CacheFile
	err = json.Unmarshal(bs, &cacheFileDB)
	if err != nil {
		return map[string]CacheFile{}, fmt.Errorf(customError.WrapError("error in reading file " + seedCacheFilePath, err))
	}
	return cacheFileDB, nil
}

func GetID(fileUrl string) string {
	hash := sha1.New()
	hash.Write([]byte(fileUrl))
	bs := hash.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

func GetFileName(respHeader http.Header, url url.URL) string {
	for headerKeys, headerValues := range respHeader {
		if headerKeys == "Content-Disposition" {
			for _, headerValue := range headerValues {
				if strings.Contains(headerValue, "filename=") {
					return strings.Split(headerValue, "filename=")[1]
				}
			}
		}
	}
	segments := strings.Split(url.Path, "/")
	if len(segments) == 0 {
		return "index"
	}
	return segments[len(segments)-1]
}

func GetCacheFileStatus(fileStatus FileStatus) string {
	switch fileStatus {
	case 0:
		return "INITIALISED"
	case 1:
		return "DOWNLOADING"
	case 2:
		return "ERROR IN DOWLOADING"
	case 3:
		return "DOWNLOADED"
	default:
		return ""
	}
}