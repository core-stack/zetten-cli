package config

import (
	"os"
	"path/filepath"
)

var home, _ = os.UserHomeDir()
var DEFAULT_ROOT_PATH = filepath.Join(home, ".zetten")
var DEFAULT_ROOT_CONFIG_PATH = filepath.Join(DEFAULT_ROOT_PATH, "config.yml")
var DEFAULT_ROOT_PACKAGES_PATH = filepath.Join(DEFAULT_ROOT_PATH, "packages")
var DEFAULT_AUTH_FILE_NAME = "zetten-auth.yml"
var DEFAULT_ROOT_AUTHS_PATH = filepath.Join(DEFAULT_ROOT_PATH, DEFAULT_AUTH_FILE_NAME)
