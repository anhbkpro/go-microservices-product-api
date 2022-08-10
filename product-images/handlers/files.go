package handlers

import (
	"github.com/anhbkpro/go-microservices-product-api/product-images/files"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
)

type Files struct {
	log   hclog.Logger
	store files.Storage
}

// NewFiles creates a new File handler
func NewFiles(s files.Storage, l hclog.Logger) *Files {
	return &Files{store: s, log: l}
}

// UploadREST implements the http.Handler interface
func (f *Files) UploadREST(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fn := vars["filename"]

	f.log.Info("Handle POST", "id", id, "filename", fn)

	// check that the filepath is a valid name and file
	if id == "" || fn == "" {
		f.invalidURI(r.URL.String(), rw)
		return
	}

	f.saveFile(id, fn, rw, r.Body)
}

func (f *Files) UploadMultipart(rw http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(128 * 1024)
	if err != nil {
		f.log.Error("Bad request", "error", err)
		http.Error(rw, "Expected multipart form data", http.StatusBadRequest)
		return
	}

	id, idErr := strconv.Atoi(r.FormValue("id"))
	f.log.Info("Process form for id", "id", id)

	if idErr != nil {
		f.log.Error("Bad request", "error", err)
		http.Error(rw, "Expected integer id", http.StatusBadRequest)
		return
	}

	ff, mh, err := r.FormFile("file")
	if err != nil {
		f.log.Error("Bad request", "error", err)
		http.Error(rw, "Expected file", http.StatusBadRequest)
		return
	}

	f.saveFile(string(id), mh.Filename, rw, ff)

}

func (f *Files) invalidURI(uri string, rw http.ResponseWriter) {
	f.log.Error("Invalid path", "path", uri)
	http.Error(rw, "Invalid file path should be in the format: /{id}/{filename}", http.StatusBadRequest)
}

func (f *Files) saveFile(id, path string, rw http.ResponseWriter, r io.ReadCloser) {
	f.log.Info("Save file for product", "id", id, "path", path)

	fp := filepath.Join(id, path)
	err := f.store.Save(fp, r)
	if err != nil {
		f.log.Error("Unable to save files", "error", err)
		http.Error(rw, "Unable to save file", http.StatusInternalServerError)
	}
}
