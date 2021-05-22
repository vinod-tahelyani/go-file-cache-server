package controllers

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"go-file-cache-server.example.com/db"
	customError "go-file-cache-server.example.com/error"
	"go-file-cache-server.example.com/models"
)

const (
	NO_AUTH      = "No Auth"
	BASIC_AUTH   = "Basic Auth"
	BEARER_TOKEN = "Bearer Token"
)

const (
	DOWNLOADS_FILE_DIR = "./downloads"
)

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func CacheFile(file models.CacheFile) {
	if file.Status == models.DOWNLOADED {
		return
	}
	file.Status = models.DOWNLOADING
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, file.FileURL, nil)
	if err != nil {
		file.Status = models.ERROR_DOWLOADING
		log.Println(customError.WrapError("error in downloading the file from '"+req.URL.String()+"'", err))
		return
	}
	setAuthHeaders(req, file)
	resp, err := client.Do(req)
	if err != nil {
		file.Status = models.ERROR_DOWLOADING
		log.Println(customError.WrapError("error in downloading the file from '"+req.URL.String()+"'", err))
		return
	}
	file.FileName = models.GetFileName(resp.Header, *req.URL)
	file.LocalPath = DOWNLOADS_FILE_DIR + "/" + file.FileName
	downloadedFile, err := os.Create(file.LocalPath)
	if err != nil {
		log.Println(customError.WrapError("error in downloading the file '"+file.FileName+"'", err))
		return
	}
	defer resp.Body.Close()
	fmt.Println("copying file")
	n, err := io.Copy(downloadedFile, resp.Body)
	if err != nil {
		log.Println(customError.WrapError("error in storing the file '"+file.FileName+"'", err))
		return
	}
	log.Printf("successfuly downloaded the file url: %s, filename: %s, size: %d", req.URL.String(), file.FileName, n)
	file.Status = models.DOWNLOADED
	db.UpdateCacheFileInDB(file)
}

func setAuthHeaders(req *http.Request, file models.CacheFile) {
	req.Header.Add("Authorization", "Basic "+basicAuth(file.Username, file.Password))
}
