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

	realPath1 := filepath.Join(baseDir, filePath)
	if fileExist(realPath1) {
		return realPath1
	}

	realPath2, err := filepath.Abs(filePath)
	if err != nil {
		panic(err)
	}
	if fileExist(realPath2) {
		return realPath2
	}

	return realPath1
}
