// // package handlers

// // import (
// // 	// "net/http"

// // 	"github.com/gorilla/mux"
// // )

// // // SetupRoutes initializes all API routes
// // func SetupRoutes() *mux.Router {
// // 	router := mux.NewRouter()

// // 	// Routes
// // 	router.HandleFunc("/upload", UploadHandler).Methods("POST")
// // 	router.HandleFunc("/download", DownloadHandler).Methods("GET")
// // 	router.HandleFunc("/delete", DeleteHandler).Methods("DELETE")

// // 	return router
// // }

// package handlers

// import (
// 	"file_mgmt_system/internal/storage"

// 	"github.com/gorilla/mux"
// 	"github.com/jmoiron/sqlx"
// )

// func NewRouter(storage storage.Storage, db *sqlx.DB) *mux.Router {
// 	router := mux.NewRouter()

// 	uploadHandler := NewUploadHandler(storage)
// 	downloadHandler := NewDownloadHandler(storage)
// 	deleteHandler := NewDeleteHandler(storage)
// 	loginHandler := NewLoginHander()

// 	router.Handle("/upload", uploadHandler).Methods("POST")
// 	router.Handle("/download/{name}", downloadHandler).Methods("GET")
// 	router.Handle("/delete/{name}", deleteHandler).Methods("DELETE")
// 	router.Handle("/storePIDetails", loginHandler).Methods("POST")

//		return router
//	}
package handlers
