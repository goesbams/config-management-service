package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/goesbams/config-management-service/models"
	"github.com/goesbams/config-management-service/utils"
)

// in-memory storing
var configStore = make(map[string]models.Config)

func CreateConfig(w http.ResponseWriter, r *http.Request) {
	log.Println("request create configuration")

	// read input
	var newConfig models.Config
	err := json.NewDecoder(r.Body).Decode(&newConfig)
	if err != nil {
		log.Println("error decoding request body:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validate config
	if err := utils.ValidateConfig(newConfig, configStore); err != nil {
		log.Println("invalid configuration data")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// store config with a version
	newConfig.Version = 1
	configStore[newConfig.Name] = newConfig

	// config created response
	log.Printf("new config created: %s", newConfig.Name)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newConfig)
}

func UpdateConfig(w http.ResponseWriter, r *http.Request) {

}

func RollbackConfig(w http.ResponseWriter, r *http.Request) {

}

func FetchConfig(w http.ResponseWriter, r *http.Request) {

}

func ListVersionsHandler(w http.ResponseWriter, r *http.Request) {

}
