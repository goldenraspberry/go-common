package config

import (
	"os"
	"path/filepath"
)

func fileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func toAbsFile(baseDir string, filePath string) string {
	if filepath.IsAbs(filePath) {
		return filePath
	}

	realPath := filepath.Join(baseDir, filePath)
	if fileExist(realPath) {
		return realPath
	}

	return realPath
}
