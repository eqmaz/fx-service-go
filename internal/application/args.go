package application

import (
	"flag"
	c "fx-service/pkg/console"
	"fx-service/pkg/e"
	util "fx-service/pkg/helpers"
	"path/filepath"
)

func getConfigFileFromArgs() string {
	// Define and parse the flags
	configPath := flag.String("config", "", "Path to the config file")
	flag.Parse()

	if *configPath == "" {
		return ""
	} else {
		return *configPath
	}
}

// tryConfigFile checks if the config file path exists, and the file is readable
// Will resolve the full path and return true if the file exists and is readable
func tryConfigFile(configFilePath string) (string, bool) {
	// Resolve the path and check if the file exists
	fullPath, err := filepath.Abs(configFilePath)
	if err != nil {
		return "", false
	}

	// Check if the config file is readable
	if !util.IsReadableFile(fullPath) {
		c.Warnf("Config file is not readable (check permissions): %s", fullPath)
		return "", false
	}

	return fullPath, true
}

// getConfigFilePath attempts to locate the config file in the current working directory and the executable directory
// First checks if a CLI argument was provided, then tries the default paths
// Config file name defaults to "config.json"
func getConfigFilePath(name string) (string, error) {
	if name == "" {
		name = "config.json"
	}

	// Get the config file path from the command line arguments
	configPath := getConfigFileFromArgs()
	if configPath != "" {
		fullPath, ok := tryConfigFile(configPath)
		if ok {
			c.Successf("Using config file from argument at: '%s'", fullPath)
			return fullPath, nil
		} else {
			c.Warnf("No config file at '%s'...", configPath)
		}
	}

	// Default paths to try the config file, if not provided in the arguments
	// These are relative to the current working directory!
	pathsToTry := []string{
		"./" + name,
		"./config/" + name,
	}

	// Check if the executable is in the current working directory already
	// If this fails, we will not care, and just try the default paths
	exeInCwd, err := util.IsExeInCwd()
	if err == nil && !exeInCwd {
		// If not (cwd is not the same as the exe dir), try the exe dir for the config files as well
		exeDir, err := util.ExeDir()
		if err == nil {
			pathsToTry = append(pathsToTry, filepath.Join(exeDir, name))
			pathsToTry = append(pathsToTry, filepath.Join(exeDir, "config/"+name))
		}
	}

	// Try the possible config paths, starting with the provided path
	for _, path := range pathsToTry {
		c.Infof("Trying config file found at: '%s'", path)
		fullPath, ok := tryConfigFile(path)
		if ok {
			return fullPath, nil
		}
	}

	return "", e.FromCode("eNcF01").SetField("pathsTried", pathsToTry)
}
