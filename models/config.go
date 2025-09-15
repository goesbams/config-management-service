package models

type Config struct {
	Name    string                 `json:"name"`
	Data    map[string]interface{} `json:"data"`
	Version int                    `json:"version"`
}
