package utils

import (
	"errors"

	"github.com/goesbams/config-management-service/models"
)

func ValidateConfig(config models.Config, configStore map[string]models.Config) error {
	if config.Name == "" {
		return errors.New("config name is required")
	}

	if config.Data == nil {
		return errors.New("config data is required")
	}

	if _, found := configStore[config.Name]; found {
		return errors.New("config already exists")
	}

	return nil
}
