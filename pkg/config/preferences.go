package config

// Internal settings which define behaviour of the Config package
const (
	NoErrorOnMissingFile = false // When true, missing config file is not an error
	//ImportUnknownFields  = false // When true, unknown fields in the config file are imported into the Config struct
)

// Priority order for configuration sources
// The lower the number, the higher the precedence. For example source 1 will overwrite values from source 2.
const (
	PriorityEnv  = 1
	PriorityFile = 2
)
