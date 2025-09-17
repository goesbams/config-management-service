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

func setupConfig() models.Config {
	return models.Config{
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
}

func TestUpdateConfig(t *testing.T) {
	handlers.ConfigStore = make(map[string]models.Config)
	handlers.ConfigStore["Main Database Config"] = setupConfig()

	t.Run("success update config", func(t *testing.T) {
		update := models.Config{
			Name: "Main Database Config",
			Type: "DATABASE",
			Versions: []models.ConfigVersion{
				{
					Version: 2,
					Property: map[string]interface{}{
						"host": "127.0.0.1",
					},
				},
			},
		}
		data, _ := json.Marshal(update)
		req, _ := http.NewRequest("POST", "/config/update", bytes.NewBuffer(data))
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handlers.UpdateConfig)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var updatedConfig models.Config
		_ = json.NewDecoder(rr.Body).Decode(&updatedConfig)
		assert.Equal(t, 2, updatedConfig.Versions[0].Version)
	})

	t.Run("error when config not found", func(t *testing.T) {
		update := models.Config{
			Name: "Unknown Config",
			Type: "DATABASE",
			Versions: []models.ConfigVersion{
				{Version: 1, Property: map[string]interface{}{"host": "127.0.0.1"}},
			},
		}
		data, _ := json.Marshal(update)
		req, _ := http.NewRequest("POST", "/config/update", bytes.NewBuffer(data))
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handlers.UpdateConfig)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Equal(t, "config not found", strings.TrimSpace(rr.Body.String()))
	})
}

func TestRollbackConfig(t *testing.T) {
	handlers.ConfigStore = make(map[string]models.Config)
	handlers.ConfigStore["Main Database Config"] = setupConfig()
	update := models.Config{
		Name: "Main Database Config",
		Type: "DATABASE",
		Versions: []models.ConfigVersion{
			{Version: 2, Property: map[string]interface{}{"host": "127.0.0.1"}},
		},
	}
	handlers.ConfigStore["Main Database Config"] = update

	t.Run("success rollback", func(t *testing.T) {
		reqBody := `{"name":"Main Database Config","version":1}`
		req, _ := http.NewRequest("POST", "/config/rollback", bytes.NewBuffer([]byte(reqBody)))
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handlers.RollbackConfig)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var rolledBack models.Config
		_ = json.NewDecoder(rr.Body).Decode(&rolledBack)
		assert.Equal(t, 3, rolledBack.Versions[len(rolledBack.Versions)-1].Version)
	})

	t.Run("error invalid version", func(t *testing.T) {
		reqBody := `{"name":"Main Database Config","version":5}`
		req, _ := http.NewRequest("POST", "/config/rollback", bytes.NewBuffer([]byte(reqBody)))
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handlers.RollbackConfig)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "invalid version", strings.TrimSpace(rr.Body.String()))
	})
}

func TestFetchConfig(t *testing.T) {
	handlers.ConfigStore = make(map[string]models.Config)
	handlers.ConfigStore["Main Database Config"] = setupConfig()

	t.Run("fetch latest version", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/config/fetch?name=Main%20Database%20Config", nil)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handlers.FetchConfig)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var fetched models.ConfigVersion
		_ = json.NewDecoder(rr.Body).Decode(&fetched)
		assert.Equal(t, 1, fetched.Version)
	})

	t.Run("fetch specific version", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/config/fetch?name=Main%20Database%20Config&version=1", nil)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handlers.FetchConfig)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var fetched models.ConfigVersion
		_ = json.NewDecoder(rr.Body).Decode(&fetched)
		assert.Equal(t, 1, fetched.Version)
	})

	t.Run("error config not found", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/config/fetch?name=Unknown", nil)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handlers.FetchConfig)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Equal(t, "config not found", strings.TrimSpace(rr.Body.String()))
	})

	t.Run("error version not found", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/config/fetch?name=Main%20Database%20Config&version=5", nil)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handlers.FetchConfig)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Equal(t, "version not found", strings.TrimSpace(rr.Body.String()))
	})
}

func TestListVersionsHandler(t *testing.T) {
	handlers.ConfigStore = make(map[string]models.Config)
	handlers.ConfigStore["Main Database Config"] = setupConfig()

	t.Run("list all versions", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/config/versions?name=Main%20Database%20Config", nil)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handlers.ListVersionsHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var versions []models.ConfigVersion
		_ = json.NewDecoder(rr.Body).Decode(&versions)
		assert.Len(t, versions, 1)
		assert.Equal(t, 1, versions[0].Version)
	})

	t.Run("error config not found", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/config/versions?name=Unknown", nil)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handlers.ListVersionsHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "config not found", strings.TrimSpace(rr.Body.String()))
	})
}
