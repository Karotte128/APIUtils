package config

import (
	"os"

	"github.com/pelletier/go-toml/v2"
)

// Loads the config from a .yaml file.
func ReadConfigFromFile(path string) (error, map[string]any) {
	// Reading the provided file path
	data, err := os.ReadFile(path)
	if err != nil {
		return err, nil
	}

	// Unmarshal the file contents into a map[string]any
	var config map[string]any
	err = toml.Unmarshal(data, &config)
	if err != nil {
		return err, nil
	}

	// no error, yay!
	return nil, config
}

// Replaces ${VAR} or ${VAR:-default} in the config with environment variables.
func ExpandEnvConfig(input map[string]any) map[string]any {
	return expandEnvRecursive(input).(map[string]any)
}
