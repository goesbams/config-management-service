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
	log.Println("request create config")

	// read input
	var newConfig models.Config
	err := json.NewDecoder(r.Body).Decode(&newConfig)
	if err != nil {
		log.Println("error decoding request body:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validate config
	if err := utils.ValidateConfig(newConfig); err != nil {
		log.Println("invalid config data", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// store new config
	configStore[newConfig.Name] = newConfig

	log.Printf("new config created: %s", newConfig.Name)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newConfig)
}

func UpdateConfig(w http.ResponseWriter, r *http.Request) {
	log.Println("request update config")

	var updatedConfig models.Config
	err := json.NewDecoder(r.Body).Decode(&updatedConfig)
	if err != nil {
		log.Println("error decoding request body:", err)
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if err := utils.ValidateConfig(updatedConfig); err != nil {
		log.Println("invalid config data")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	existingConfig, found := configStore[updatedConfig.Name]
	if !found {
		log.Println("config not found")
		http.Error(w, "config not found", http.StatusNotFound)
		return
	}

	// increment latest version
	ev := existingConfig.Versions
	newVersion := ev[len(ev)-1].Version + 1
	updatedConfig.Versions = append(updatedConfig.Versions,
		models.ConfigVersion{
			Version:  newVersion,
			Property: updatedConfig.Versions[0].Property,
		},
	)

	// store updated config
	configStore[updatedConfig.Name] = updatedConfig

	log.Printf("config updated: %s, version: %d", updatedConfig.Name, newVersion)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedConfig)
}

func RollbackConfig(w http.ResponseWriter, r *http.Request) {
	log.Println("request rollback config")

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

	// retrieve existing config
	existingConfig, found := configStore[rollbackRequest.Name]
	if !found {
		log.Println("config not found")
		http.Error(w, "config not found", http.StatusBadRequest)
		return
	}

	// check if version exists & is less than current version
	ev := existingConfig.Versions
	if rollbackRequest.Version <= 0 || rollbackRequest.Version >= ev[len(ev)-1].Version {
		log.Println("invalid version request")
		http.Error(w, "invalid version", http.StatusBadRequest)
		return
	}

	// get specific version
	rollbackVersion := ev[rollbackRequest.Version-1]

	// create new version
	newVersion := ev[len(ev)-1].Version + 1

	// update rollback config
	rollbackConfig := existingConfig
	rollbackConfig.Versions = append(rollbackConfig.Versions,
		models.ConfigVersion{
			Version:  newVersion,
			Property: rollbackVersion.Property,
		},
	)

	log.Printf("config rollback: %s, version: %d", rollbackConfig.Name, newVersion)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(rollbackConfig)
}

func FetchConfig(w http.ResponseWriter, r *http.Request) {
	log.Println("request fetch config")

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
