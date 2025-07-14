package util

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

type CreateOptions struct {
	Path        string
	IsDir       bool
	Perm        os.FileMode
	InitialData string
}

type CreateOpt func(*CreateOptions)

func WithPath(path string) CreateOpt {
	return func(o *CreateOptions) {
		o.Path = path
	}
}

func AsDir() CreateOpt {
	return func(o *CreateOptions) {
		o.IsDir = true
	}
}

func WithPerm(perm os.FileMode) CreateOpt {
	return func(o *CreateOptions) {
		o.Perm = perm
	}
}

func WithInitialData(data string) CreateOpt {
	return func(o *CreateOptions) {
		o.InitialData = data
	}
}

func CreateIfNotExistsByPath(opts ...CreateOpt) error {
	options := &CreateOptions{}
	for _, opt := range opts {
		opt(options)
	}

	if _, err := os.Stat(options.Path); os.IsNotExist(err) {
		if options.IsDir {
			fmt.Println(fmt.Printf("📁 %s not found, creating...", options.Path))
			if err := os.MkdirAll(options.Path, options.Perm); err != nil {
				return fmt.Errorf("❌ Failed to create %s: %w", options.Path, err)
			}
		} else {
			fmt.Println(fmt.Printf("📝 %s not found, creating...", options.Path))
			if err := os.WriteFile(options.Path, []byte(options.InitialData), options.Perm); err != nil {
				return fmt.Errorf("❌ Failed to create %s: %w", options.Path, err)
			}
		}
	}
	return nil
}

func SaveYAMLIndented(path string, data any) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := yaml.NewEncoder(f, yaml.Indent(4))
	defer encoder.Close()

	return encoder.Encode(data)
}
