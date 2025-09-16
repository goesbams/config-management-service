package models

type ConfigType string

const (
	databaseConfig    ConfigType = "DATABASE"
	apiKeysConfig     ConfigType = "API_KEY"
	loggingConfig     ConfigType = "LOGGING"
	featureFlagConfig ConfigType = "FEATURE_FLAG"
	networkSetting    ConfigType = "NETWORK_SETTING"
)

func (c ConfigType) IsValid() bool {
	switch c {
	case databaseConfig, apiKeysConfig, loggingConfig, featureFlagConfig, networkSetting:
		return true
	}

	return false
}

type Config struct {
	Name     string          `json:"name"`
	Type     ConfigType      `type:"type"`
	Versions []ConfigVersion `json:"versions"`
}

type ConfigVersion struct {
	Version  int                    `json:"version"`
	Property map[string]interface{} `json:"property"`
}
