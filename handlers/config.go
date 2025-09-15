package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/goesbams/config-management-service/models"
)

// in-memory storing
var configStore = make(map[string]models.Config)

func CreateConfig(w http.ResponseWriter, r *http.Request) {
	// read input
	var newConfig models.Config
	err := json.NewDecoder(r.Body).Decode(&newConfig)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validate config
	if newConfig.Name == "" || newConfig.Data == nil {
		http.Error(w, "Invalid config", http.StatusBadRequest)
		return
	}

	// store config with a version
	newConfig.Version = 1
	configStore[newConfig.Name] = newConfig

	// config created response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newConfig)
}
