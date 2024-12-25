package handlers

import (
	"file_mgmt_system/internal/storage"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type DeleteHandler struct {
	Storage storage.Storage
	// Producer *kafka.KafkaProducer
}

func NewDeleteHandler(s storage.Storage) *DeleteHandler {
	return &DeleteHandler{
		Storage: s,
		// Producer: producer,
	}
}

func (h *DeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	vars := mux.Vars(r)
	// objectName := r.URL.Query().Get("name")
	objectName := vars["name"]
	if objectName == "" {
		http.Error(w, "File name is required", http.StatusBadRequest)
		return
	}

	err := h.Storage.DeleteFile(objectName)
	if err != nil {
		http.Error(w, "Failed to delete file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// message := map[string]string{
	// 	"event":   "file_deleted",
	// 	"file_id": objectName,
	// }
	// messageBytes, _ := json.Marshal(message)

	// err = h.Producer.SendMessage(objectName, string(messageBytes))
	// if err != nil {
	// 	http.Error(w, "Failed to send Kafka message", http.StatusInternalServerError)
	// 	return
	// }

	// fmt.Fprintf(w, "File %s deleted successfully", objectName)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("File %s deleted successfully", objectName)))
}
