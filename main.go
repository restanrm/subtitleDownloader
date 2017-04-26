package main

import (
	"flag"
	"log"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/restanrm/subtitleDownloader/filefinder"
	"github.com/restanrm/subtitleDownloader/opensubtitle"
	"github.com/spf13/viper"
)

type Finder interface {
	Find(rootpath string) string
}

type SubFetcher interface {
	//FetchTo fetch filepath subtitle to destination
	FetchTo(filepath, destination string) error
}

func init() {
	viper.SetDefault("language", "en")
	viper.SetDefault("filefinder.basepath", "/data/videos/")
	viper.SetDefault("opensubtitle.ratelimit", time.Second*10/40)
}

func main() {
	verbose := flag.Bool("v", false, "Set verbosity to higher level")
	flag.Parse()
	if *verbose {
		logrus.SetLevel(logrus.DebugLevel)
	}

	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/subtitleDownloader")
	viper.AddConfigPath("$HOME/.subtitleDownloader")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Couldn't read configuration")
	}

	// create chan of file that need subtitle
	fileNeedSub := make(chan string, 100)
	go filefinder.Find(viper.GetString("filefinder.basepath"), fileNeedSub)

	if viper.IsSet("opensubtitle.username") && viper.IsSet("opensubtitle.password") {
		username := viper.GetString("opensubtitle.username")
		password := viper.GetString("opensubtitle.password")
		language := viper.GetString("language")

		ratelimiter := time.Tick(viper.GetDuration("opensubtitle.ratelimit"))
		o := ost.New(username, password, language)

		for path := range fileNeedSub {
			<-ratelimiter
			target := filefinder.GetSubFilepath(path)
			err = o.Fetch(path, target)
			if err != nil {
				ctx := logrus.Fields{
					"error":  err,
					"path":   path,
					"target": target,
				}
				logrus.WithFields(ctx).Error("Failed to retrieve subtitle")
			}
		}
	}

}
