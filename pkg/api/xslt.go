package api

import (
	"fmt"

	"mime"
	"strings"

	"net/http"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

// Tranform XML using XSLT.
// Not implemented yet.
//
func XSLT(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var path = r.URL.Path
	if path == "/xslt" {
		path = "/xslt/index.html"
	}

	if strings.HasSuffix(path, "/") {
		path = fmt.Sprintf("%sindex.html", path)
	}

	path = strings.TrimPrefix(path, "/xslt")

	mimeType := "text/plain"
	pathParts := strings.Split(path, ".")
	if len(pathParts) > 1 {
		mimeType = mime.TypeByExtension("." + pathParts[1])
	}

	bs, err := Asset(path)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}

	w.Header().Set("Content-Type", mimeType)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(bs)
	if err != nil {
		log.Error("failed to write", err)
	}

}
