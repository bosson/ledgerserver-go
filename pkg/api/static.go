package api

import (
	"errors"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"strings"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

// Static returns files stored in bindata.go
func Static(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var path = r.URL.Path
	if path == "/web/" {
		path = "/web/index.html"
	}
	if strings.HasSuffix(path, "/") {
		path = fmt.Sprintf("%sindex.html", path)
	}

	path = strings.TrimPrefix(path, "/web/")

	mimeType := "text/plain"
	pathParts := strings.Split(path, ".")
	if len(pathParts) > 1 {
		mimeType = mime.TypeByExtension("." + pathParts[1])
	}

	bs, err := Asset(path)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", mimeType)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(bs)
	if err != nil {
		log.Error("failed to write", err)
	}

}

// Asset used to read asset files
func Asset(f string) (b []byte, err error) {

	log.Info("load asset: ", f)

	if strings.Contains(f, "..") {
		return nil, errors.New("path with '..' not allowed")
	}

	file, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return ioutil.ReadAll(file)

}
