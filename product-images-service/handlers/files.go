package handlers

import (
	"io"
	"net/http"
	"path/filepath"
	"product-images-service/files"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
)

type Files struct {
	logger hclog.Logger
	store  files.Storage
}

func NewFilesHandler(logger hclog.Logger, store files.Storage) *Files {
	return &Files{
		logger,
		store,
	}
}

// Upload REST implementation
func (f *Files) UploadREST(rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	filename := vars["filename"]

	f.logger.Info("Handle POST", "id", id, "filename", filename)

	f.saveFile(id, filename, rw, r.Body)
}

// Upload multipart file
func (f *Files) UploadMultipart(rw http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(128 * 1024)

	if err != nil {
		f.logger.Error("Bad request", err)
		http.Error(rw, "Expected multipart form data", http.StatusBadRequest)
	}

	id, idErr := strconv.Atoi(r.FormValue("id"))

	f.logger.Info("Process form for id", "id", id)
	file, mh, err := r.FormFile("file")

	if idErr != nil {
		f.logger.Error("Bad request", err)
		http.Error(rw, "Expected integer Id", http.StatusBadRequest)
	}
	if err != nil || idErr != nil {
		f.logger.Error("Bad request", err)
		http.Error(rw, "Expected file", http.StatusBadRequest)
	}

	f.saveFile(r.FormValue("id"), mh.Filename, rw, file)
}

func (f *Files) saveFile(
	id string,
	filename string,
	rw http.ResponseWriter,
	r io.ReadCloser,
) {
	f.logger.Info("Saving file", "id", id, "filename", filename)

	fp := filepath.Join(id, filename)

	err := f.store.Save(fp, r)

	if err != nil {
		f.logger.Error("Unable to save file", "error", err)
		http.Error(rw, "Unable to save file", http.StatusInternalServerError)
	}
}

func (f *Files) invalidURI(uri string, rw http.ResponseWriter) {
	f.logger.Error("Invalid path", "path", uri)
	http.Error(rw, "Invalid file path should be in the format: /[id]/[filepath]", http.StatusBadRequest)
}
