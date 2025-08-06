package pkg

import (
	"fmt"

	"github.com/core-stack/zetten-cli/internal/core/file"
)

type PackageConfig struct {
	PackageFile `yaml:",inline"`
}

func LoadPackageConfig(path string) (*PackageConfig, error) {
	cfg, err := file.Load[PackageConfig](path)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%+v", cfg)

	cfg.Path = path
	return cfg, nil
}
