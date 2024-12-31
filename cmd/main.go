package main

import (
	"context"
	"file_mgmt_system/internal/handlers"
	"file_mgmt_system/internal/kafka"
	"file_mgmt_system/internal/service"
	"file_mgmt_system/internal/storage"
	"file_mgmt_system/middleware"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/oracle/oci-go-sdk/v65/common"
)

func main() {

	provider := common.DefaultConfigProvider()
	ociStorage, err := storage.NewOCIStorage(provider, "test_bucket_1")
	if err != nil {
		fmt.Println("Failed to initialize storage:", err)
		return
	}
	// ---------

	dbUser := "root"
	dbPassword := ""
	dbHost := "127.0.0.1:3306"
	dbName := "dms2"

	db, err := storage.NewDB(dbUser, dbPassword, dbHost, dbName)
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer db.Conn.Close()

	//----------

	redisAddr := "127.0.0.1:6379"
	redisPassword := ""
	redisDB := 0

	redisClient, err := storage.NewRedisClient(redisAddr, redisPassword, redisDB)
	if err != nil {
		log.Fatalf("Error initializing Redis: %v", err)
	}
	defer redisClient.Close()

	// -------------
	brokers := []string{"localhost:9092"}
	topic := "file-events"
	groupID := "file-mgmt-group"

	producer := kafka.NewKafkaProducer(brokers, topic)
	defer producer.Close()

	consumer := kafka.NewKafkaConsumer(brokers, topic, groupID)
	defer consumer.Close()

	// Start Kafka consumer in a goroutine
	go func() {
		ctx := context.Background()
		log.Println("Starting Kafka consumer...")
		consumer.ConsumeMessages(ctx)
	}()
	// ----------
	router := mux.NewRouter()

	loginService := service.NewLoginService(db)
	loginHandler := handlers.NewLoginHandler(loginService)
	router.Handle("/storePIDetails", http.HandlerFunc(loginHandler.ServeHTTP)).Methods("POST")
	router.Handle("/getEmail", http.HandlerFunc(loginHandler.GetEmail)).Methods("GET")

	uploadService := service.NewUploadService(db, ociStorage)
	uploadHandler := handlers.NewUploadHandler(uploadService)
	router.Handle("/upload", uploadHandler).Methods("POST")

	downloadService := service.NewDownloadService(db, ociStorage, redisClient)
	downloadHandler := handlers.NewDownloadHandler(downloadService)
	router.Handle("/downloadFile/{fileID}", http.HandlerFunc(downloadHandler.ServeHTTP)).Methods("GET")
	router.Handle("/getThumbnail/{fileName}", http.HandlerFunc(downloadHandler.GetThumbnail)).Methods("GET")
	router.Handle("/getPreview/{fileID}", http.HandlerFunc(downloadHandler.GetPreview)).Methods("GET")

	getFileList := service.NewGetFilesService(db)
	getFileListHandler := handlers.NewGetFilesHandler(getFileList)
	router.Handle("/getUploadedFiles/{email}", getFileListHandler).Methods("GET")

	deleteService := service.NewDeleteService(db, ociStorage)
	deleteHandler := handlers.NewDeleteHandler(deleteService)
	router.Handle("/deleteFile/{fileID}", deleteHandler).Methods("DELETE")

	// ------------
	// router := handlers.NewRouter(ociStorage, db.Conn)
	corsRouter := middleware.CORS(router)
	cookiesRouter := middleware.CookieMiddleware(corsRouter)

	// -------
	fmt.Println("Server is running on :8081")
	err = http.ListenAndServe("localhost:8081", cookiesRouter)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	// -------
}
