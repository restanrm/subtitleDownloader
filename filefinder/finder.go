package filefinder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/fsnotify/fsnotify"
)

func Find(rootpath string, outpath chan string) {

	//go watch(rootpath, outpath)
	err := filepath.Walk(rootpath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			ctx := logrus.Fields{
				"err":      err,
				"path":     path,
				"fileinfo": info,
				"rootpath": rootpath,
			}
			logrus.WithFields(ctx).Error("An error occurred while walking through directories")
		}

		ext := filepath.Ext(path)
		moviesExtensions := []string{".mkv", ".avi"}
		if !inSlice(ext, moviesExtensions) {
			return nil
		}

		if IsSrtForFile(path) {
			return nil
		}

		outpath <- path

		return nil
	})
	if err != nil {
		ctx := logrus.Fields{
			"err": err,
		}
		logrus.WithFields(ctx).Error("An error occured while retrieving some files")
	}
	close(outpath)
}

func watch(rootpath string, outpath chan string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logrus.WithField("err", err).Error("Failed to create a new fsnotify watcher")
		return
	}
	defer watcher.Close()

	done := make(chan struct{})
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				ctx := logrus.Fields{
					"rootpath": rootpath,
				}
				logrus.WithFields(ctx).Info("event received:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					fmt.Println("FileModified", event.Name)
					path := event.Name
					ext := filepath.Ext(path)
					moviesExtensions := []string{".mkv", ".avi"}
					if !inSlice(ext, moviesExtensions) {
						continue
					}

					if IsSrtForFile(path) {
						continue
					}

					//outpath <- path
				}
			case err := <-watcher.Errors:
				logrus.WithField("err", err).Error("Received an error from err channel on watcher")
			}
		}
	}()

	err = watcher.Add(rootpath)
	if err != nil {
		logrus.WithFields(logrus.Fields{"err": err, "rootpath": rootpath}).Error("Failed to add rootpath to watcher")
	}
	<-done
}

func inSlice(a string, b []string) (ret bool) {
	ctx := logrus.Fields{
		"a":     a,
		"slice": b,
	}
	logrus.WithFields(ctx).Debug("Is value in slice")
	for _, c := range b {
		if c == a {
			return true
		}
	}
	return
}

func GetSubFilepath(path string) string {
	base := filepath.Base(path)
	dir := filepath.Dir(path)
	ext := filepath.Ext(path)
	withoutExt := strings.TrimSuffix(base, ext)

	ctx := logrus.Fields{
		"base":       base,
		"dir":        dir,
		"ext":        ext,
		"withoutExt": withoutExt,
	}
	logrus.WithFields(ctx).Debug("Searching for subtitle")

	return filepath.Join(dir, withoutExt+".en.srt")
}

func IsSrtForFile(path string) bool {
	target := GetSubFilepath(path)
	fi, err := os.Stat(target)
	if err != nil {
		ctx := logrus.Fields{
			"err":    err,
			"path":   path,
			"target": target,
		}
		logrus.WithFields(ctx).Debug("File doesn't exist")
		return false
	}

	if fi.Name() == target {
		return true
	}
	return true

}
