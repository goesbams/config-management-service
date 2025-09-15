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
	log.Println("request update configuration")

	var updatedConfig models.Config
	err := json.NewDecoder(r.Body).Decode(&updatedConfig)
	if err != nil {
		log.Println("error decoding request body:", err)
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if err := utils.ValidateConfig(updatedConfig, configStore); err != nil {
		log.Println("invalid configuration data")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	existingConfig, found := configStore[updatedConfig.Name]
	if !found {
		log.Println("config not found")
		http.Error(w, "config not found", http.StatusNotFound)
		return
	}

	updatedConfig.Version = existingConfig.Version + 1
	configStore[updatedConfig.Name] = updatedConfig

	log.Printf("config updated: %s, version: %d", updatedConfig.Name, updatedConfig.Version)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedConfig)
}

func RollbackConfig(w http.ResponseWriter, r *http.Request) {
	log.Println("request for rollback configuration")

	var rollbackRequest struct {
		Name    string `json:"name"`
		Version int    `json:"version"`
	}

	err := json.NewDecoder(r.Body).Decode(&rollbackRequest)
	if err != nil {
		log.Println("error decoding request body:", err)
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	existingConfig, exists := configStore[rollbackRequest.Name]
	if !exists || rollbackRequest.Version >= existingConfig.Version {
		log.Println("invalid version")
		http.Error(w, "invalid version", http.StatusBadRequest)
		return
	}

	rollbackConfig := existingConfig
	rollbackConfig.Version = rollbackRequest.Version
	configStore[rollbackRequest.Name] = rollbackConfig

	log.Printf("config rollback: %s, version: %d", rollbackConfig.Name, rollbackConfig.Version)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(rollbackConfig)
}

func FetchConfig(w http.ResponseWriter, r *http.Request) {
	log.Println("request fetch configuration")

	configName := r.URL.Query().Get("name")
	if configName == "" {
		http.Error(w, "config name is required", http.StatusBadRequest)
		return
	}

	config, exists := configStore[configName]
	if !exists {
		log.Println(w, "config not found")
		http.Error(w, "config not found", http.StatusNotFound)
		return
	}

	log.Printf("latest config: %s", configName)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(config)
}

func ListVersionsHandler(w http.ResponseWriter, r *http.Request) {

}
