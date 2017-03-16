package config

import (
	"os"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
)

// GetLogLevel returns the log level
func GetLogLevel() log.Level {
	level := getEnv("ENM_LOG_LEVEL", "")

	switch level {
	case "Debug":
		return log.DebugLevel
	case "Info":
		return log.InfoLevel
	case "Warn":
		return log.WarnLevel
	case "Error":
		return log.ErrorLevel
	case "Fatal":
		return log.FatalLevel
	case "Panic":
		return log.PanicLevel
	}

	return log.InfoLevel
}

// GetDependentLogLevel returns the log level for dependent devices
func GetDependentLogLevel() log.Level {
	level := getEnv("DEPENDENT_LOG_LEVEL", "")

	switch level {
	case "Debug":
		return log.DebugLevel
	case "Info":
		return log.InfoLevel
	case "Warn":
		return log.WarnLevel
	case "Error":
		return log.ErrorLevel
	case "Fatal":
		return log.FatalLevel
	case "Panic":
		return log.PanicLevel
	}

	return log.InfoLevel
}

// GetLoopDelay returns the time delay in seconds between each application process loop
func GetLoopDelay() (time.Duration, error) {
	value, err := strconv.Atoi(getEnv("ENM_CONFIG_LOOP_DELAY", "10"))
	return time.Duration(value), err
}

// GetPauseDelay returns the time delay in seconds between each pause check
func GetPauseDelay() (time.Duration, error) {
	value, err := strconv.Atoi(getEnv("ENM_CONFIG_PAUSE_DELAY", "10"))
	return time.Duration(value), err
}

// GetShortBluetoothTimeout returns the timeout for each instantaneous bluetooth operation
func GetShortBluetoothTimeout() (time.Duration, error) {
	value, err := strconv.Atoi(getEnv("ENM_BLUETOOTH_SHORT_TIMEOUT", "1"))
	return time.Duration(value), err
}

// GetLongBluetoothTimeout returns the timeout for each running bluetooth operation
func GetLongBluetoothTimeout() (time.Duration, error) {
	value, err := strconv.Atoi(getEnv("ENM_BLUETOOTH_LONG_TIMEOUT", "10"))
	return time.Duration(value), err
}

// GetAssetsDir returns the root directory used to store the database and application commits
func GetAssetsDir() string {
	return getEnv("ENM_ASSETS_DIRECTORY", "/data/assets")
}

// GetDbDir returns the directory used to store the database
func GetDbDir() string {
	return getEnv("ENM_DB_DIRECTORY", "/data/database")
}

// GetDbName returns the directory used to store the database
func GetDbName() string {
	return getEnv("ENM_DB_NAME", "my.db")
}

// GetVersion returns the API version used to communicate with the supervisor
func GetVersion() string {
	return getEnv("ENM_API_VERSION", "v1")
}

// GetSuperAddr returns the address used to communicate with the supervisor
func GetSuperAddr() string {
	return getEnv("RESIN_SUPERVISOR_ADDRESS", "http://127.0.0.1:4000")
}

// GetSuperAPIKey returns the API key used to communicate with the supervisor
func GetSuperAPIKey() string {
	return getEnv("RESIN_SUPERVISOR_API_KEY", "")
}

// GetLockFileLocation returns the location of the lock file
func GetLockFileLocation() string {
	return getEnv("ENM_LOCK_FILE_LOCATION", "/data/resin-updates.lock")
}

func getEnv(key, fallback string) string {
	result := os.Getenv(key)
	if result == "" {
		result = fallback
	}
	return result
}
