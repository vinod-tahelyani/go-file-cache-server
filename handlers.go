package main

import (
	"encoding/json"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"go-file-cache-server.example.com/controllers"
	"go-file-cache-server.example.com/db"
	"go-file-cache-server.example.com/error"
	"go-file-cache-server.example.com/models"
)

func healthCheckHandler(respWriter http.ResponseWriter, req *http.Request) {
	respWriter.WriteHeader(http.StatusOK)
	respWriter.Write([]byte("OK"))
}

func cacheFile(respWriter http.ResponseWriter, req *http.Request) {
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		respBody := error.GetHTTPError("could not read request body", err)
		bs, err := json.Marshal(respBody)
		if err != nil {
			log.Println(err)
			respWriter.WriteHeader(http.StatusInternalServerError)
			return
		}
		respWriter.Write(bs)
	}
	// make new CacheFileRequestBody
	var cacheFileRequest models.CacheFileRequestBody
	json.Unmarshal(bs, &cacheFileRequest)
	cacheFileRequest.ID = models.GetID(cacheFileRequest.FileURL)
	if cacheFileRequest.AuthType == "" {
		cacheFileRequest.AuthType = controllers.NO_AUTH
	}

	// make new cacheFile
	cacheFile, err := db.NewCacheFile(cacheFileRequest)
	var cacheFileRespBody models.CacheFileResponseBody
	var respStatus int
	if err != nil {
		cacheFileRespBody = models.CacheFileResponseBody{
			Message: err.Error(),
		}
		respStatus = http.StatusBadRequest
	} else {
		db.AddCacheFileToDB(cacheFile)
		go controllers.CacheFile(cacheFile)
		cacheFileRespBody = models.CacheFileResponseBody{
			ID:      cacheFile.ID,
			Status:  cacheFile.Status,
			Message: "Initialised successfully",
		}
		respStatus = http.StatusOK
	}
	bs, _ = json.Marshal(cacheFileRespBody)
	respWriter.Header().Add("Content-Type", "application/json")
	respWriter.WriteHeader(respStatus)
	respWriter.Write(bs)
}

func getCachedFileStatus(respWriter http.ResponseWriter, req *http.Request) {
	fileURL := req.URL.Query().Get("fileURL")
	id := req.URL.Query().Get("ID")
	type responseBody struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	}
	var respBody responseBody
	var respStatus int
	if fileURL == "" && id == "" {
		respBody = responseBody{
			Message: "/cache-file-status requires fileURL or ID as query parameter",
		}
		respStatus = http.StatusBadRequest
	} else if len(id) > 0 {
		cacheFile, ok := db.GetCacheFileByID(id)
		if !ok {
			respBody = responseBody{
				Message: "could not find file by id " + id,
			}
			respStatus = http.StatusBadRequest
		} else {
			respBody = responseBody{
				Status: string(models.GetCacheFileStatus(cacheFile.Status)),
			}
			respStatus = http.StatusOK
		}
	} else {
		if cacheFile, ok := db.GetCacheFileByURL(fileURL); !ok {
			respBody = responseBody{
				Message: "could not find file by fileURL " + fileURL,
			}
			respStatus = http.StatusBadRequest
		} else {
			respBody = responseBody{
				Status: string(models.GetCacheFileStatus(cacheFile.Status)),
			}
			respStatus = http.StatusOK
		}
	}
	bs, _ := json.Marshal(respBody)
	respWriter.Header().Add("Content-Type", "application/json")
	respWriter.WriteHeader(respStatus)
	respWriter.Write(bs)
}

func getFile(respWriter http.ResponseWriter, req *http.Request) {
	fileURL := req.URL.Query().Get("fileURL")
	id := req.URL.Query().Get("ID")
	type responseBody struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	}
	var respBody responseBody
	var respStatus int
	var filePath string = ""
	var fileName string = ""
	if fileURL == "" && id == "" {
		respBody = responseBody{
			Message: "/cache-file-status requires fileURL or ID as query parameter",
		}
		respStatus = http.StatusBadRequest
	} else if len(id) > 0 {
		cacheFile, ok := db.GetCacheFileByID(id)
		if !ok {
			respBody = responseBody{
				Message: "could not find file by id " + id,
			}
			respStatus = http.StatusBadRequest
		} else {
			filePath = cacheFile.LocalPath
			fileName = cacheFile.FileName
			respStatus = http.StatusOK
		}
	} else {
		if cacheFile, ok := db.GetCacheFileByURL(fileURL); !ok {
			respBody = responseBody{
				Message: "could not find file by fileURL " + fileURL,
			}
			respStatus = http.StatusBadRequest
		} else {
			filePath = cacheFile.LocalPath
			fileName = cacheFile.FileName
			respStatus = http.StatusOK
		}
	}
	if len(filePath) == 0 {
		bs, _ := json.Marshal(respBody)
		respWriter.Header().Add("Content-Type", "application/json")
		respWriter.WriteHeader(respStatus)
		respWriter.Write(bs)
	} else {
		respWriter.Header().Add("Content-Type", "application/json")
		respWriter.Header().Set("Content-Disposition", "attachment; filename="+fileName)
		respWriter.Header().Set("content-type", "application/octet-stream")
		respWriter.WriteHeader(respStatus)
		file, err := os.Open(filePath)
		if err != nil {
			log.Fatal("error in opening file: ", err)
		}
		io.Copy(respWriter, file)
	}
}

func invalidateCache(respWriter http.ResponseWriter, req *http.Request) {
	os.RemoveAll(controllers.DOWNLOADS_FILE_DIR)
	os.Mkdir(controllers.DOWNLOADS_FILE_DIR, fs.ModeDevice)
	os.Chmod(controllers.DOWNLOADS_FILE_DIR, os.ModePerm)
	db.InitializeDB()
	respWriter.Header().Add("Content-Type", "application/json")
	respWriter.WriteHeader(http.StatusOK)
	type HTTPResponse struct {
		Message string `json:"message"`
	}
	respBody := HTTPResponse {
		Message: "OK",
	}
	bs, _ := json.Marshal(respBody)
	respWriter.Write(bs)
}
