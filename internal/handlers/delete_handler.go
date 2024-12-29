package handlers

import (
	"encoding/json"
	"file_mgmt_system/internal/service"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type DeleteHandler struct {
	Service *service.DeleteService
	// Producer *kafka.KafkaProducer
}

func NewDeleteHandler(service *service.DeleteService) *DeleteHandler {
	return &DeleteHandler{
		Service: service,
		// Producer: producer,
	}
}

type DeleteResponse struct {
	Success  bool   `json:"success"`
	FileName string `json:"fileName,omitempty"`
	Message  string `json:"message,omitempty"`
}

func (h *DeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(DeleteResponse{
			Success: false,
			Message: "Invalid request method",
		})
		return
	}

	vars := mux.Vars(r)
	fileID, ok := vars["fileID"]
	if !ok || fileID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(DeleteResponse{
			Success: false,
			Message: "fileID is required in the path",
		})
		return
	}

	fileName, err := h.Service.Delete(fileID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(DeleteResponse{
			Success: false,
			Message: "Failed to delete file: " + err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(DeleteResponse{
		Success:  true,
		FileName: fileName,
		Message:  fmt.Sprintf("File %s deleted successfully", fileName),
	})
}
