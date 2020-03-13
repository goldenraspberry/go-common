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
	cfg Config

	mapDefaultConfig = map[string]string{
		"env":        "release",
		"listen":     ":8080",
		"cache_dir":  "/tmp/cache",
		"tmp_dir":    "/tmp",
		"log_dir":    "./logs",
		"log_level":  "INFO",
		"log":        "app.log",
		"access_log": STDOUT,
		"error_log":  STDERR,
		"slow_log":   DISABLE,
	}
)

func GetListen() string {
	return cfg.listen
}

func GetEnv() string {
	return cfg.env
}

func GetConfig(section string) map[string]string {
	lock.RLock()
	defer lock.RUnlock()

	v, _ := cfg.config.GetSection(section)
	if v == nil {
		return map[string]string{}
	}
	return v
}

func GetCacheDir() string {
	lock.RLock()
	defer lock.RUnlock()
	return cfg.cacheDir
}

func GetTmpDir() string {
	lock.RLock()
	defer lock.RUnlock()
	return cfg.tmpDir
}

func GetLogDir() string {
	lock.RLock()
	defer lock.RUnlock()
	return cfg.logDir
}

func GetLogPath() string {
	lock.RLock()
	defer lock.RUnlock()
	return cfg.logPath
}

func GetAccessLogPath() string {
	lock.RLock()
	defer lock.RUnlock()
	return cfg.accessLogPath
}

func GetErrorLogPath() string {
	lock.RLock()
	defer lock.RUnlock()
	return cfg.errorLogPath
}

func GetSlowLogPath() string {
	lock.RLock()
	defer lock.RUnlock()
	return cfg.slowLogPath
}

func GetLogLevel() string {
	lock.RLock()
	defer lock.RUnlock()
	return cfg.logLevel
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

	loadConfig(goConfig)
}

func loadConfig(goConfig *goconfig.ConfigFile) {
	cfg.config = goConfig

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
	cfg.env = mapGlobal["env"]
	cfg.listen = mapGlobal["listen"]
	cfg.cacheDir = toAbsFile(baseDir, mapGlobal["cache_dir"])
	cfg.tmpDir = toAbsFile(baseDir, mapGlobal["tmp_dir"])
	cfg.logDir = toAbsFile(baseDir, mapGlobal["log_dir"])
	cfg.logLevel = mapGlobal["log_level"]

	cfg.logPath = getLogFile(cfg.logDir, mapGlobal["log"])
	cfg.accessLogPath = getLogFile(cfg.logDir, mapGlobal["access_log"])
	cfg.errorLogPath = getLogFile(cfg.logDir, mapGlobal["error_log"])
	cfg.slowLogPath = getLogFile(cfg.logDir, mapGlobal["slow_log"])
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
