package prune_test

import (
	"prunehild/prune"
	"reflect"
	"testing"
	"testing/fstest"
	"time"
)

func must[T any](obj T, err error) T {
	if err != nil {
		panic(err)
	}
	return obj
}

func makeTime(s string) time.Time {
	return must(time.Parse(time.DateTime, s))
}

func TestPrune(t *testing.T) {
	// Given
	fs := fstest.MapFS{
		"1": &fstest.MapFile{
			ModTime: makeTime("2009-11-01 10:00:00"),
		},
		"2": &fstest.MapFile{
			ModTime: makeTime("2009-11-01 11:00:00"),
		},
		"3": &fstest.MapFile{
			ModTime: makeTime("2009-11-01 12:00:00"),
		},
		"4": &fstest.MapFile{
			ModTime: makeTime("2009-11-01 13:00:00"),
		},
		"5": &fstest.MapFile{
			ModTime: makeTime("2009-11-01 14:00:00"),
		},
		"6": &fstest.MapFile{
			ModTime: makeTime("2009-11-01 15:00:00"),
		},
		"7": &fstest.MapFile{
			ModTime: makeTime("2009-11-01 16:00:00"),
		},
	}
	var calledFiles []string

	// When
	prune.Prune(fs, must(fs.Glob("*")), func(fn string) {
		calledFiles = append(calledFiles, fn)
	})

	expectedFilenames := []string{"3", "4"}
	if !reflect.DeepEqual(expectedFilenames, calledFiles) {
		t.Errorf("got %v, want %v", calledFiles, expectedFilenames)
	}
}
