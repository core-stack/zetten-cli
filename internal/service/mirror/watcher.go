package mirror

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/core-stack/zetten-cli/internal/core/root"
	"github.com/fsnotify/fsnotify"
)

type Config struct {
	Mirror [][]string `yaml:"mirror"`
}

var (
	wg      sync.WaitGroup
	stopCh  = make(chan struct{})
	watched = make(map[string]struct{})
)

func StartMirrorService(configPath string) error {
	cfg, err := root.LoadRootConfig()
	if err != nil {
		return err
	}
	for _, group := range cfg.Mirror {
		wg.Add(1)
		go watchGroup(group)
	}

	wg.Wait()
	return nil
}

func StopMirrorService() {
	close(stopCh)
	wg.Wait()
}

func watchGroup(paths []string) {
	defer wg.Done()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	for _, root := range paths {
		addRecursive(watcher, root)
	}

	log.Printf("Monitoring paths: %v", paths)

	for {
		select {
		case <-stopCh:
			return
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			source := findSourceRoot(paths, event.Name)
			if source == "" {
				continue
			}
			relPath, _ := filepath.Rel(source, event.Name)

			for _, destRoot := range paths {
				if destRoot == source {
					continue
				}
				destPath := filepath.Join(destRoot, relPath)

				switch {
				case event.Op&fsnotify.Create == fsnotify.Create:
					fi, err := os.Stat(event.Name)
					if err == nil && fi.IsDir() {
						addRecursive(watcher, event.Name)
					} else {
						copyFile(event.Name, destPath)
					}
				case event.Op&fsnotify.Write == fsnotify.Write:
					copyFile(event.Name, destPath)
				case event.Op&fsnotify.Remove == fsnotify.Remove:
					os.RemoveAll(destPath)
				case event.Op&fsnotify.Rename == fsnotify.Rename:
					os.RemoveAll(destPath)
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("Erro watcher:", err)
		}
	}
}

func findSourceRoot(roots []string, path string) string {
	for _, root := range roots {
		if strings.HasPrefix(path, root) {
			return root
		}
	}
	return ""
}

func addRecursive(watcher *fsnotify.Watcher, path string) {
	filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil || !info.IsDir() {
			return nil
		}
		if _, ok := watched[p]; ok {
			return nil
		}
		err = watcher.Add(p)
		if err != nil {
			log.Printf("Error adding path %s: %v", p, err)
		} else {
			watched[p] = struct{}{}
		}
		return nil
	})
}

func copyFile(src, dst string) {
	err := os.MkdirAll(filepath.Dir(dst), 0755)
	if err != nil {
		log.Printf("Erro criando diretório: %v", err)
		return
	}
	srcFile, err := os.Open(src)
	if err != nil {
		return
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		log.Printf("Erro criando destino: %v", err)
		return
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		log.Printf("Erro copiando %s → %s: %v", src, dst, err)
	}
}
