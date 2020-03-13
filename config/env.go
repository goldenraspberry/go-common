package config

import (
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Env struct {
	configFile string
	env        string
	listen     string
	logDir     string
	tmpDir     string
	baseDir    string
	cacheDir   string
}

var (
	env = &Env{}
)

func envEnv() string {
	return env.env
}

func envListen() string {
	return env.listen
}

func GetBaseDir() string {
	return envBaseDir()
}

func envBaseDir() string {
	return env.baseDir
}

func envTmpDir() string {
	return env.tmpDir
}

func envCacheDir() string {
	return env.cacheDir
}

func envLogDir() string {
	return env.logDir
}

func envConfigFile() string {
	return env.configFile
}

func initEnv() {
	e := &Env{}

	defaultBaseDir := getDefaultBasePath()

	appEnv := flag.String("env", "", "app env, default release")
	listen := flag.String("listen", "", "web server [ip]:port, default :8080")
	configFile := flag.String("config-file", "", "app base config, default env.ini")
	baseDir := flag.String("base-dir", defaultBaseDir, "app base dir")
	tmpDir := flag.String("tmp-dir", "", "app tmp dir, default tmp")
	cacheDir := flag.String("cache-dir", "", "app base dir, default cache")
	logDir := flag.String("log-dir", "", "app base dir, default, log")

	flag.Parse()

	e.env = *appEnv
	e.listen = *listen
	e.baseDir = *baseDir
	e.tmpDir = getRelPathForEnv(*tmpDir)
	e.cacheDir = getRelPathForEnv(*cacheDir)
	e.logDir = getRelPathForEnv(*logDir)
	e.configFile = getRelPathForEnv(*configFile)

	env = e
}

func getBinaryDir() string {
	curFilename := os.Args[0]
	binaryPath, err := exec.LookPath(curFilename)
	if err != nil {
		panic(err)
	}

	binaryPath, err = filepath.Abs(binaryPath)
	if err != nil {
		panic(err)
	}

	basePath := filepath.Dir(binaryPath)

	return basePath
}

func getDefaultBasePath() string {
	basePath := getBinaryDir()

	// if binary file is in "bin" dir, get the parent dir
	if strings.HasSuffix(basePath, "/bin") {
		basePath = filepath.Dir(basePath)
	}

	return basePath
}

func getRelPathForEnv(filePath string) string {
	if filePath == "" {
		return ""
	}

	if filepath.IsAbs(filePath) {
		return filePath
	}

	realPath, err := filepath.Abs(filePath)
	if err != nil {
		panic(err)
	}

	return realPath
}
