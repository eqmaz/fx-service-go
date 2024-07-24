package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// Private functions
// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -

// mergeConfig - merges source (config) into target and overwrites target fields, where needed
func mergeConfig(target, source map[string]interface{}) map[string]interface{} {
	for key, value := range source {
		target[key] = value
	}
	return target
}

// setDefaults - sets the default values from defaultConfigMap to a Config map
func setDefaults(config map[string]interface{}) {
	for key, value := range defaultConfigMap {
		config[key] = value
	}
}

// convertMapToStruct - converts a map to a struct
func convertMapToStruct(m map[string]interface{}, result interface{}) error {
	bytes, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, result)
}

// NewConfig - returns a new loaded, merged config
func NewConfig(filePath string) (Config, error) {
	defaults := make(map[string]interface{})
	setDefaults(defaults)

	var fileConfig map[string]interface{}
	var envConfig map[string]interface{}
	var resultConfig = defaults

	var err error

	// Load from the config file if a file is specified
	if filePath != "" {
		fileConfig, err = LoadConfigFromFile(filePath)
		if err != nil {
			fmt.Printf("Error loading config from file: %v. Using defaults.\n", err)
		} else {
			resultConfig = mergeConfig(resultConfig, fileConfig)
		}
	}

	// Load environment variables
	envConfig = LoadConfigFromEnv()

	// Merge configurations based on priority
	if PriorityEnv < PriorityFile {
		resultConfig = mergeConfig(resultConfig, envConfig)
		if filePath != "" {
			resultConfig = mergeConfig(resultConfig, fileConfig)
		}
	} else {
		if filePath != "" {
			resultConfig = mergeConfig(resultConfig, fileConfig)
		}
		resultConfig = mergeConfig(resultConfig, envConfig)
	}

	// Convert the final map to a Config struct
	var finalConfig Config
	err = convertMapToStruct(resultConfig, &finalConfig)
	if err != nil {
		return Config{}, err
	}

	return finalConfig, nil
}

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// Public functions
// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -

// DefaultConfig - returns the default configuration as a Config map
func DefaultConfig() map[string]interface{} {
	config := make(map[string]interface{})
	setDefaults(config)
	return config
}

// LoadConfigFromFile - Load config from file
func LoadConfigFromFile(filePath string) (map[string]interface{}, error) {
	file, err := os.Open(filePath)
	if err != nil {
		if NoErrorOnMissingFile {
			return nil, nil
		}
		return nil, fmt.Errorf("could not open config file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("could not close config file: %v\n", err)
		}
	}(file)

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("could not read config file: %v", err)
	}

	var fileConfig map[string]interface{}
	if err := json.Unmarshal(bytes, &fileConfig); err != nil {
		return nil, fmt.Errorf("could not parse config file: %v", err)
	}

	return fileConfig, nil
}

// LoadConfigFromEnv - Load config from environment variables
// Only loads environment variables that were defined in the defaultConfigMap
func LoadConfigFromEnv() map[string]interface{} {
	envConfig := make(map[string]interface{})

	for key := range defaultConfigMap {
		if value, exists := os.LookupEnv(key); exists {
			envConfig[key] = value
		}
	}

	return envConfig
}
