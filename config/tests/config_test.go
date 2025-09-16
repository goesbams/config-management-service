package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/goesbams/config-management-service/handlers"
	"github.com/goesbams/config-management-service/models"
)

var configStore = make(map[string]models.Config)

func TestCreateConfig(t *testing.T) {
	t.Run("success create new config", func(t *testing.T) {
		config := models.Config{
			Name: "Main Database Config",
			Type: "DATABASE",
			Versions: []models.ConfigVersion{
				{
					Version: 1,
					Property: map[string]interface{}{
						"host":     "localhost",
						"port":     5432,
						"username": "admin",
						"password": "secret",
					},
				},
			},
		}

		configData, _ := json.Marshal(config)

		req, err := http.NewRequest("POST", "/config", bytes.NewBuffer(configData))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handlers.CreateConfig)

		handler.ServeHTTP(rr, req)

		var actualConfig models.Config
		err = json.NewDecoder(rr.Body).Decode(&actualConfig)
		assert.NoError(t, err)
		assert.Equal(t, config.Name, actualConfig.Name)
		assert.Equal(t, config.Type, actualConfig.Type)
	})

	t.Run("error when config name is empty", func(t *testing.T) {
		config := models.Config{
			Name: "",
			Type: "DATABASE",
			Versions: []models.ConfigVersion{
				{
					Version: 1,
					Property: map[string]interface{}{
						"host":     "localhost",
						"port":     5432,
						"username": "admin",
						"password": "secret",
					},
				},
			},
		}

		configData, _ := json.Marshal(config)

		req, err := http.NewRequest("POST", "/config", bytes.NewBuffer(configData))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handlers.CreateConfig)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "config name is required", strings.TrimSpace(rr.Body.String()))
	})

	t.Run("error when config type is empty", func(t *testing.T) {
		config := models.Config{
			Name: "Main Database Config",
		}

		configData, _ := json.Marshal(config)

		req, err := http.NewRequest("POST", "/config", bytes.NewBuffer(configData))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handlers.CreateConfig)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "config type is required", strings.TrimSpace(rr.Body.String()))
	})

	t.Run("error when config type is invalid", func(t *testing.T) {
		config := models.Config{
			Name: "Main Database Config",
			Type: "OTHERS",
		}

		configData, _ := json.Marshal(config)

		req, err := http.NewRequest("POST", "/config", bytes.NewBuffer(configData))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handlers.CreateConfig)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		fmt.Println(rr.Body.String())
		assert.Equal(t, "invalid config type", strings.TrimSpace(rr.Body.String()))
	})

	t.Run("error when versions is nil", func(t *testing.T) {
		config := models.Config{
			Name: "Main Database Config",
			Type: "DATABASE",
		}

		configData, _ := json.Marshal(config)

		req, err := http.NewRequest("POST", "/config", bytes.NewBuffer(configData))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handlers.CreateConfig)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "invalid config versions")
	})

	t.Run("Error: Config versions length is 0", func(t *testing.T) {
		config := models.Config{
			Name:     "Main Database Config",
			Type:     "DATABASE",
			Versions: []models.ConfigVersion{},
		}

		configData, _ := json.Marshal(config)

		req, err := http.NewRequest("POST", "/config", bytes.NewBuffer(configData))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handlers.CreateConfig)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "config versions is required")
	})
}

func TestUpdateConfig(t *testing.T) {

	t.Run("success update config and increment version", func(t *testing.T) {
		config := models.Config{
			Name: "Main Database Config",
			Type: "DATABASE",
			Versions: []models.ConfigVersion{
				{
					Version: 1,
					Property: map[string]interface{}{
						"host":     "localhost",
						"port":     5432,
						"username": "admin",
						"password": "secret",
					},
				},
			},
		}

		configStore[config.Name] = config

		updatedConfig := models.Config{
			Name: "Main Database Config",
			Type: "DATABASE",
			Versions: []models.ConfigVersion{
				{
					Property: map[string]interface{}{
						"host":     "localhost",
						"port":     5432,
						"username": "admin",
						"password": "new_secret",
					},
				},
			},
		}

		updatedData, _ := json.Marshal(updatedConfig)

		req, err := http.NewRequest("PUT", "/config/update", bytes.NewBuffer(updatedData))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handlers.UpdateConfig)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		var actualConfig models.Config
		err = json.NewDecoder(rr.Body).Decode(&actualConfig)
		assert.NoError(t, err)

		assert.Len(t, actualConfig.Versions, 2)
		assert.Equal(t, 2, actualConfig.Versions[1].Version)
	})

	t.Run("error when config not found", func(t *testing.T) {
		updatedConfig := models.Config{
			Name: "Non Existing Config",
			Type: "DATABASE",
			Versions: []models.ConfigVersion{
				{
					Property: map[string]interface{}{
						"host":     "localhost",
						"port":     5432,
						"username": "admin",
						"password": "new_secret",
					},
				},
			},
		}

		updatedData, _ := json.Marshal(updatedConfig)

		req, err := http.NewRequest("PUT", "/config/update", bytes.NewBuffer(updatedData))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handlers.UpdateConfig)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Contains(t, rr.Body.String(), "config not found")
	})
}
