package ost

import (
	"testing"

	"github.com/spf13/viper"
)

func TestFetch(t *testing.T) {
	o := New(viper.GetString("opensubtitle.username"), viper.GetString("opensubtitle.password"), "eng")
	if o == nil {
		t.Error("Failed to get ost instance quit!")
	}

	path := "/data/videos/series/Arrow/s05/Arrow.S05E18.Disbanded_SDTV.mkv"
	err := o.Fetch(path)
	if err != nil {
		t.Error("Failed to retrieve subtitle for path: %v", path)
	}

}
