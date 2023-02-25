package config

import (
	"fmt"
	"github.com/arielsrv/go-archaius"
	"os"
	"path"
	"runtime"

	"github.com/src/main/app/log"

	"github.com/src/main/app/helpers/files"

	"github.com/src/main/app/config/env"
)

const (
	File = "config.yml"
)

func init() {
	showWd()
	log.Warn("trying to load config ...")
	_, caller, _, _ := runtime.Caller(0)
	root := path.Join(path.Dir(caller), "../../..")
	err := os.Chdir(root)
	if err != nil {
		showWd()
		wd, wdErr := os.Getwd()
		if wdErr != nil {
			log.Fatal(wdErr)
		}
		root = path.Join(wd, "/src")
	}

	propertiesPath, environment, scope :=
		fmt.Sprintf("%s/resources/config", root),
		env.GetEnv(),
		env.GetScope()

	var compositeConfig []string

	scopeConfig := fmt.Sprintf("%s/%s/%s.%s", propertiesPath, environment, scope, File)
	if files.Exist(scopeConfig) {
		compositeConfig = append(compositeConfig, scopeConfig)
	}

	envConfig := fmt.Sprintf("%s/%s/%s", propertiesPath, environment, File)
	if files.Exist(envConfig) {
		compositeConfig = append(compositeConfig, envConfig)
	}

	sharedConfig := fmt.Sprintf("%s/%s", propertiesPath, File)
	if files.Exist(sharedConfig) {
		compositeConfig = append(compositeConfig, sharedConfig)
	}

	err = archaius.Init(
		archaius.WithENVSource(),
		archaius.WithRequiredFiles(compositeConfig),
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Infof("INFO: ENV: %s, SCOPE: %s", environment, scope)
}

func showWd() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	log.Info("working directory: " + wd)
}

func String(key string) string {
	value, err := archaius.GetValue(key).ToString()
	if err != nil {
		fallback := ""
		log.Warnf("warn: config %s not found, fallback to empty string", key)
		return fallback
	}
	return value
}

func Int(key string) int {
	value, err := archaius.GetValue(key).ToInt()
	if err != nil {
		var fallback = 0
		log.Warnf(fmt.Sprintf("warn: config %s not found, fallback to %d", key, fallback))
		return fallback
	}
	return value
}

func TryInt(key string, defaultValue int) int {
	value, err := archaius.GetValue(key).ToInt()
	if err != nil {
		log.Warnf(fmt.Sprintf("warn: config %s not found, fallback to %d", key, defaultValue))
		return defaultValue
	}
	return value
}
