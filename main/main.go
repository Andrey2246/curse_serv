package main

import (
	"log"
	"net/http"
	"os"

	"curse_serv/handlers"
)

func main() {
	// Ensure the storage directory exists
	const storageDir = "/home/cursework/fileStorage"
	err := os.MkdirAll(storageDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create storage directory: %v", err)
	}

	// Register HTTP routes
	http.HandleFunc("/upload", handlers.UploadHandler(storageDir))
	http.HandleFunc("/list", handlers.ListHandler)
	http.HandleFunc("/download", handlers.DownloadHandler(storageDir))
	http.HandleFunc("/wipe", handlers.WipeHandler(storageDir))

	defer handlers.WipeHandler(storageDir)

	// Start the server
	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
