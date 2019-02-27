package api

import (
	"fmt"
	"os"
	"path"

	"github.com/bosson/ledgerserver"
	"github.com/julienschmidt/httprouter"

	"net/http"
)

// Version returns the servers version
func Version(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// http.StatusOK, "text/plain"
	fmt.Fprintf(w, "%s/%s", path.Base(os.Args[0]), ledgerserver.Version)

}
