package logging

import (
	"errors"
	"strings"

	"github.com/sirupsen/logrus"
)

var log = NewMasterLogger()

const (
	// logModuleKey is the key used for the module name data entry
	logModuleKey = "_module"
	// logPriorityKey is the log entry key for priority log statements
	logPriorityKey = "_priority"
	// logPriorityCritical is the log entry value for priority log statements
	logPriorityCritical = "CRITICAL"
)

// LevelFromString returns a logrus.Level from a string identifier
func LevelFromString(s string) (logrus.Level, error) {
	switch strings.ToLower(s) {
	case "debug":
		return logrus.DebugLevel, nil
	case "info", "notice":
		return logrus.InfoLevel, nil
	case "warn", "warning":
		return logrus.WarnLevel, nil
	case "error":
		return logrus.ErrorLevel, nil
	case "fatal", "critical":
		return logrus.FatalLevel, nil
	case "panic":
		return logrus.PanicLevel, nil
	default:
		return logrus.DebugLevel, errors.New("could not convert string to log level")
	}
}

var (
	globalPackageLevels []PkgLogConfig
)

// MustGetLogger returns a package-aware logger from the master logger
func MustGetLogger(module string) *Logger {
	if globalPackageLevels == nil {
		var err error
		if globalPackageLevels, err = parsePkgLogLevelsFromEnv(); err == nil {
			SetPackageLevels(globalPackageLevels)
		}
	}
	return log.PackageLogger(module)
}

// AddHook adds a hook to the global logger
func AddHook(hook logrus.Hook) {
	log.AddHook(hook)
}

// EnableColors enables colored logging
func EnableColors() {
	log.EnableColors()
}

// DisableColors disables colored logging
func DisableColors() {
	log.DisableColors()
}

// SetLevel sets the logger's minimum log level
func SetLevel(level logrus.Level) {
	log.SetLevel(level)
}

// SetPackageLevels sets the independent log levels for packages
func SetPackageLevels(config []PkgLogConfig) {
	log.Debugf("Setting up logging config %v", config)
	ConfigPkgLogging(log, config)
	log.Debugf("Logging package levels %v", log.PkgConfig)
}
