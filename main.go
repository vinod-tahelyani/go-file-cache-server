package main

import (
	"io/fs"
	"net/http"
	"os"

	"go-file-cache-server.example.com/controllers"
	"go-file-cache-server.example.com/db"
)

func main()  {
	db.InitializeDB()
	os.Mkdir(controllers.DOWNLOADS_FILE_DIR, fs.ModeDevice)
	os.Chmod(controllers.DOWNLOADS_FILE_DIR, os.ModePerm)
	s := NewServer()
	http.ListenAndServe(s.addr, s)
}

func (s server) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	requestURI := request.URL.Path
	for _, route := range routes {
		matches := route.regex.FindStringSubmatch(requestURI)
		if len(matches) > 0 && route.method == request.Method {
			route.handler(w, request)
			return
		}
	}
	http.NotFound(w, request)
}

type server struct {
	addr string
}

func NewServer() server {
	return server{addr: ":8080"}
}