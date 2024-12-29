package handlers

import (
	"file_mgmt_system/internal/service"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

type DownloadHandler struct {
	Service *service.DownloadService
}

func NewDownloadHandler(service *service.DownloadService) *DownloadHandler {
	return &DownloadHandler{
		Service: service,
	}
}

func (h *DownloadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	fileID, ok := vars["fileID"]
	if !ok || fileID == "" {
		http.Error(w, "fileID is required in the path", http.StatusBadRequest)
		return
	}

	fileName, fileStream, err := h.Service.Download(fileID)
	if err != nil {
		http.Error(w, "Failed to download file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// defer fileStream.Close()

	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	// w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	// w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)
	// io.Copy(w, reader)
	_, err = io.Copy(w, fileStream)
	if err != nil {
		http.Error(w, "Failed to send file: "+err.Error(), http.StatusInternalServerError)
	}
}
