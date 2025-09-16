package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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
	updatedVersion := existingConfig.Versions[len(existingConfig.Versions)-1].Version + 1
	updatedVersionConfig := models.ConfigVersion{
		Version:  updatedVersion,
		Property: updatedConfig.Versions[0].Property,
	}

	existingConfig.Versions = append(existingConfig.Versions, updatedVersionConfig)

	// store updated config
	configStore[updatedConfig.Name] = existingConfig

	log.Printf("config updated: %s, version: %d", updatedConfig.Name, updatedVersion)
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
	rollbackConfig.Versions = append(rollbackConfig.Versions, models.ConfigVersion{
		Version:  newVersion,
		Property: rollbackVersion.Property,
	})

	// Store the rollback config
	configStore[rollbackConfig.Name] = rollbackConfig

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

	// check if a spesific version requested
	versionStr := r.URL.Query().Get("version")
	if versionStr != "" {

		// convert to int
		version, err := strconv.Atoi(versionStr)
		if err != nil {
			log.Println("invalid version format")
			http.Error(w, "invalid version format", http.StatusBadRequest)
			return
		}

		if version <= 0 || version > len(config.Versions) {
			log.Println("version not found")
			http.Error(w, "version not found", http.StatusBadRequest)
			return
		}

		// specific version
		specificVersion := config.Versions[version-1]

		log.Printf("spesific config version %d: %s", version, config.Name)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(specificVersion)
		return
	}

	// latest version
	latestVersion := config.Versions[len(config.Versions)-1]

	log.Printf("latest config version %s", config.Name)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(latestVersion)
}

func ListVersionsHandler(w http.ResponseWriter, r *http.Request) {
}
