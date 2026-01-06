package configx

import (
	"os"
	"path"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func fileExist(path string) bool {
	_, err := os.Lstat(path)

	return !os.IsNotExist(err)
}

func searchConfigFile(fileName string, configPaths []string) string {
	for _, dir := range configPaths {
		if f := path.Join(dir, fileName); fileExist(f) {
			return f
		}
	}

	return ""
}

func LoadOnConfigPath(configName string, configPaths []string, cfg interface{}) (configFileUsed string, err error) {
	defaultConfigPaths, _ := GetDefaultAppConfigRoot()

	allConfigPaths := make([]string, 0, len(configPaths)+len(defaultConfigPaths))
	allConfigPaths = append(allConfigPaths, configPaths...)
	allConfigPaths = append(allConfigPaths, defaultConfigPaths...)
	for _, configPath := range configPaths {
		for _, defaultConfigPath := range defaultConfigPaths {
			allConfigPaths = append(allConfigPaths, filepath.Join(filepath.Dir(defaultConfigPath), configPath))
		}
	}

	for _, configPath := range configPaths {
		if path.IsAbs(configPath) {
			continue
		}

		for _, defaultConfigPath := range defaultConfigPaths {
			allConfigPaths = append(allConfigPaths, path.Join(defaultConfigPath, configPath))
		}
	}

	configFileUsed = searchConfigFile(configName, allConfigPaths)
	if configFileUsed == "" {
		return
	}

	f, err := os.Open(configFileUsed)
	if err != nil {
		return
	}
	defer func() {
		_ = f.Close()
	}()

	err = yaml.NewDecoder(f).Decode(cfg)

	return
}

func Load(configName string, cfg interface{}) (configFileUsed string, err error) {
	return LoadOnConfigPath(configName, []string{
		"./", "./config/", "../", "../config/", "../../", "../../config/",
	}, cfg)
}
