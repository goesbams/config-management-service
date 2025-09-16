package utils

import (
	"errors"

	"github.com/goesbams/config-management-service/models"
)

func ValidateConfig(config models.Config) error {
	if config.Name == "" {
		return errors.New("config name is required")
	}

	if config.Type == "" {
		return errors.New("config type is required")
	}

	if !config.Type.IsValid() {
		return errors.New("invalid config type")
	}

	if config.Versions == nil {
		return errors.New("invalid config versions")
	}

	if len(config.Versions) == 0 {
		return errors.New("config versions is required")
	}

	return nil
}
