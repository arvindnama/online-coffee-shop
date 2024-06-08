package handlers

import (
	"net/http"
	"path/filepath"
	"product-images-service/files"

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

func (f *Files) ServeHttp(rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	filename := vars["filename"]

	f.logger.Info("Handle POST", "id", id, "filename", filename)

	f.saveFile(id, filename, rw, r)
}

func (f *Files) saveFile(
	id string,
	filename string,
	rw http.ResponseWriter,
	r *http.Request,
) {
	f.logger.Info("Saving file", "id", id, "filename", filename)

	fp := filepath.Join(id, filename)

	err := f.store.Save(fp, r.Body)

	if err != nil {
		f.logger.Error("Unable to save file", "error", err)
		http.Error(rw, "Unable to save file", http.StatusInternalServerError)
	}
}

func (f *Files) invalidURI(uri string, rw http.ResponseWriter) {
	f.logger.Error("Invalid path", "path", uri)
	http.Error(rw, "Invalid file path should be in the format: /[id]/[filepath]", http.StatusBadRequest)
}
