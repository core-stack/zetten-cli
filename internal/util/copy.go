package util

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// CopyDir recursively copies a directory tree from source to target, skipping ignored names.
func CopyDir(source, target string, ignore []string) error {
	entries, err := os.ReadDir(source)
	if err != nil {
		return fmt.Errorf("error reading source directory %s: %w", source, err)
	}

	if err := os.MkdirAll(target, 0755); err != nil {
		return fmt.Errorf("error creating destination directory %s: %w", target, err)
	}

	for _, entry := range entries {
		name := entry.Name()

		if isIgnored(name, ignore) {
			continue
		}

		srcPath := filepath.Join(source, name)
		dstPath := filepath.Join(target, name)

		info, err := entry.Info()
		if err != nil {
			return fmt.Errorf("error getting info for %s: %w", srcPath, err)
		}

		if info.IsDir() {
			if err := CopyDir(srcPath, dstPath, ignore); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return fmt.Errorf("error copying file from %s to %s: %w", srcPath, dstPath, err)
			}
		}
	}

	return nil
}

func isIgnored(name string, ignore []string) bool {
	for _, s := range ignore {
		if s == name {
			return true
		}
	}
	return false
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}

	if fi, err := os.Stat(src); err == nil {
		return os.Chmod(dst, fi.Mode())
	}

	return nil
}
