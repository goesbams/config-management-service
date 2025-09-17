package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/goesbams/config-management-service/handlers"
)

func main() {
	// routes Setup
	http.HandleFunc("/config", handlers.CreateConfig)
	http.HandleFunc("/config/update", handlers.UpdateConfig)
	http.HandleFunc("/config/rollback", handlers.RollbackConfig)
	http.HandleFunc("/config/fetch", handlers.FetchConfig)
	http.HandleFunc("/config/versions", handlers.ListVersionsHandler)

	fmt.Println("Server started at :8090")
	log.Fatal(http.ListenAndServe(":8090", nil))
}
