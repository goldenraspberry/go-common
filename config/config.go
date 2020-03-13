package config

import (
	"strings"

	"github.com/Unknwon/goconfig"
)

const (
	STDOUT  = "stdout"
	STDERR  = "stderr"
	DISABLE = "disable"
)

type Config struct {
	listen string

	env      string
	cacheDir string
	tmpDir   string

	logDir        string
	logPath       string
	accessLogPath string
	errorLogPath  string
	slowLogPath   string
	logLevel      string

	config *goconfig.ConfigFile
}

var (
	config Config

	mapDefaultConfig = map[string]string{
		"env":        "release",
		"listen":     ":8080",
		"cache_dir":  "./cache",
		"tmp_dir":    "./tmp",
		"log_dir":    "./logs",
		"log_level":  "INFO",
		"log":        "app.log",
		"access_log": STDOUT,
		"error_log":  STDERR,
		"slow_log":   DISABLE,
	}
)

func GetListen() string {
	return config.listen
}

func GetEnv() string {
	return config.env
}

func GetConfig(section string) map[string]string {
	lock.RLock()
	defer lock.RUnlock()

	v, _ := config.config.GetSection(section)
	if v == nil {
		return map[string]string{}
	}
	return v
}

func GetCacheDir() string {
	lock.RLock()
	defer lock.RUnlock()
	return config.cacheDir
}

func GetTmpDir() string {
	lock.RLock()
	defer lock.RUnlock()
	return config.tmpDir
}

func GetLogDir() string {
	lock.RLock()
	defer lock.RUnlock()
	return config.logDir
}

func GetLogPath() string {
	lock.RLock()
	defer lock.RUnlock()
	return config.logPath
}

func GetAccessLogPath() string {
	lock.RLock()
	defer lock.RUnlock()
	return config.accessLogPath
}

func GetErrorLogPath() string {
	lock.RLock()
	defer lock.RUnlock()
	return config.errorLogPath
}

func GetSlowLogPath() string {
	lock.RLock()
	defer lock.RUnlock()
	return config.slowLogPath
}

func GetLogLevel() string {
	lock.RLock()
	defer lock.RUnlock()
	return config.logLevel
}

func initConfig() {
	lock.Lock()
	defer lock.Unlock()

	baseDir := envBaseDir()
	configFile := envConfigFile()

	if configFile == "" {
		configFile = "env.ini"
	}

	goConfig := loadConfigFile(baseDir, configFile)

	config.config = goConfig

	loadConfig()
}

func loadConfig() {
	goConfig := config.config

	global, _ := goConfig.GetSection("global")

	mapGlobal := map[string]string{
		"env":        envEnv(),
		"listen":     envListen(),
		"cache_dir":  envCacheDir(),
		"tmp_dir":    envTmpDir(),
		"log_dir":    envLogDir(),
		"log_level":  "",
		"log":        "",
		"access_log": "",
		"error_log":  "",
		"slow_log":   "",
	}

	for k := range mapGlobal {
		if _, ok := global[k]; ok {
			mapGlobal[k] = global[k]
		}
	}
	for k, v := range mapGlobal {
		if v == "" {
			if _, ok := mapDefaultConfig[k]; ok {
				mapGlobal[k] = mapDefaultConfig[k]
			}
		}
	}

	baseDir := envBaseDir()
	config.env = mapGlobal["env"]
	config.listen = mapGlobal["listen"]
	config.cacheDir = toAbsFile(baseDir, mapGlobal["cache_dir"])
	config.tmpDir = toAbsFile(baseDir, mapGlobal["tmp_dir"])
	config.logDir = toAbsFile(baseDir, mapGlobal["log_dir"])
	config.logLevel = mapGlobal["log_level"]

	config.logPath = getLogFile(config.logDir, mapGlobal["log"])
	config.accessLogPath = getLogFile(config.logDir, mapGlobal["access_log"])
	config.errorLogPath = getLogFile(config.logDir, mapGlobal["error_log"])
	config.slowLogPath = getLogFile(config.logDir, mapGlobal["slow_log"])
}

func getLogFile(baseDir, filePath string) string {
	lowerFilePath := strings.ToLower(filePath)
	switch lowerFilePath {
	case DISABLE:
		fallthrough
	case STDERR:
		fallthrough
	case STDOUT:
		return lowerFilePath
	}
	return toAbsFile(baseDir, filePath)
}
