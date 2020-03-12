package config

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Unknwon/goconfig"
	"github.com/goldenraspberry/go-common/utils"
)

func loadConfigFile(baseDir string, file string) *goconfig.ConfigFile {
	configFilePath := toAbsFile(baseDir, file)

	if fileExist(configFilePath) {
		panic(fmt.Sprintf("config file %s not exists", configFilePath))
	}

	goConfig, err := goconfig.LoadConfigFile(configFilePath)
	if err != nil {
		panic(err)
	}

	if err = loadIncludeFiles(baseDir, goConfig); err != nil {
		panic("load include files error:" + err.Error())
	}

	utils.Go(func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGUSR1)

		for {
			sig := <-ch
			switch sig {
			case syscall.SIGUSR1:
				newGoConfig := reloadConfigFile(baseDir, configFilePath)
				if newGoConfig != nil {
					lock.Lock()
					goConfig = newGoConfig
					loadConfig()
					lock.Unlock()
					publishReloadSignal()
				}
			}
		}
	})
	return goConfig
}

func reloadConfigFile(baseDir string, configFilePath string) *goconfig.ConfigFile {
	var err error
	goConfig, err := goconfig.LoadConfigFile(configFilePath)
	if err != nil {
		log.Println("reload config file, error:", err)
		return nil
	}

	if err = loadIncludeFiles(baseDir, goConfig); err != nil {
		log.Println("reload files include files error:", err)
		return nil
	}
	log.Println("reload config file successfully！")
	return goConfig
}

func loadIncludeFiles(baseDir string, goConfig *goconfig.ConfigFile) error {
	includeFile := goConfig.MustValue("include_path", "path", "")
	if includeFile != "" {
		includeFiles := strings.Split(includeFile, ",")

		incFiles := make([]string, len(includeFiles))
		for i, incFile := range includeFiles {
			configFilePath := toAbsFile(baseDir, incFile)
			if fileExist(configFilePath) {
				incFiles[i] = configFilePath
			}
		}
		return goConfig.AppendFiles(incFiles...)
	}

	return nil
}
