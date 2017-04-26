//Opensubtitle downloader
package ost

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/oz/osdb"
	"github.com/restanrm/subtitleDownloader/filefinder"
)

type Ost struct {
	*osdb.Client
	language string
}

func New(username, password, language string) *Ost {
	c, err := osdb.NewClient()
	if err != nil {
		logrus.WithField("err", err).Errorln("Failed to create New osdbClient")
		return nil
	}

	err = c.LogIn(username, password, language)
	if err != nil {
		logrus.WithField("err", err).Error("Couldn't set username and password")
		return nil
	}

	return &Ost{Client: c, language: language}
}

func (o *Ost) Fetch(path, destination string) error {
	if filefinder.IsSrtForFile(path) {
		return errors.New("File already exist")
	}
	languages := []string{o.language}
	res, err := o.Client.FileSearch(path, languages)
	if err != nil {
		return err
	}
	if len(res) < 1 {
		return fmt.Errorf("No subtitles found for %v", filepath.Base(path))
	}
	err = o.Client.DownloadTo(&res[0], destination)
	if err != nil {
		return err
	}
	logrus.Infof("Successfully downloaded %v", filepath.Base(destination))
	return nil
}
