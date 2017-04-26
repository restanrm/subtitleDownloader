package filefinder

import "testing"

func TestFinder(t *testing.T) {
	outpaths := make(chan string, 100)
	go Find("/data/videos", outpaths)
	for out := range outpaths {
		t.Log(out)
	}
}
