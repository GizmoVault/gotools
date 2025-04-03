package configx

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	ErrInvalidPath = errors.New("invalid path")
)

func GetWorkDirectory() (dir string, err error) {
	return os.Getwd()
}

func GetWorkDirectoryKey(workDir string) (key string, err error) {
	return GetWorkDirectoryKeyEx(workDir, nil)
}

func GetWorkDirectoryKeyEx(workDir string, trimPrefix []string) (key string, err error) {
	if workDir == "" {
		err = ErrInvalidPath

		return
	}

	key = workDir

	for _, prefix := range trimPrefix {
		if strings.HasPrefix(key, prefix) {
			key = key[len(prefix):]

			return
		}
	}

	if strings.HasPrefix(workDir, "/") {
		key = workDir[1:]

		return
	}

	if len(workDir) <= 2 && workDir[1] != ':' {
		err = ErrInvalidPath

		return
	}

	key = workDir[0:1] + workDir[2:]

	return
}

type ZjzConfig struct {
	ConfigRoot       string   `yaml:"config_root"`
	LocalSuffix      string   `yaml:"local_suffix"`
	TrimAppKeyPrefix []string `yaml:"trim_app_key_prefix"`
}

func GetDefaultAppConfigRoot() (roots []string, err error) {
	globalCfgRoot, err := os.UserConfigDir()
	if err != nil {
		return
	}

	d, err := os.ReadFile(path.Join(globalCfgRoot, "zjz-config.yaml"))
	if err != nil {
		return
	}

	var cfg ZjzConfig
	err = yaml.Unmarshal(d, &cfg)
	if err != nil {
		return
	}

	wd, err := GetWorkDirectory()
	if err != nil {
		return
	}

	appKey, err := GetWorkDirectoryKeyEx(wd, cfg.TrimAppKeyPrefix)
	if err != nil {
		return
	}

	root1 := filepath.Join(cfg.ConfigRoot, appKey)

	if cfg.LocalSuffix != "" {
		roots = append(roots, filepath.Join(root1, cfg.LocalSuffix))
	}

	roots = append(roots, root1)

	return
}
