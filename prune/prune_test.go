package prune_test

import (
	"prunehild/prune"
	"prunehild/schedules"
	"reflect"
	"slices"
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
	fs := fstest.MapFS{
		"01": &fstest.MapFile{ModTime: makeTime("2014-05-17 16:38:51")},
		"02": &fstest.MapFile{ModTime: makeTime("2014-05-17 16:40:11")},
		"03": &fstest.MapFile{ModTime: makeTime("2014-05-17 16:40:51")},
		"04": &fstest.MapFile{ModTime: makeTime("2014-05-17 16:41:11")},
		"05": &fstest.MapFile{ModTime: makeTime("2014-05-17 16:41:31")},
		"06": &fstest.MapFile{ModTime: makeTime("2014-05-17 16:41:46")},
		"07": &fstest.MapFile{ModTime: makeTime("2014-05-17 16:41:51")},
		"08": &fstest.MapFile{ModTime: makeTime("2014-05-17 16:41:56")},
		"09": &fstest.MapFile{ModTime: makeTime("2014-05-17 16:42:01")},
	}
	allFiles := must(fs.Glob("*"))
	slices.Sort(allFiles)
	testSched := []time.Duration{
		schedules.Second * 5,
		schedules.Second * 20,
		schedules.Second * 40,
		schedules.Second * 80,
		schedules.Long,
	}
	start := makeTime("2014-05-17 16:42:02")

	tests := []struct {
		now      time.Time
		expected []string
	}{
		{start, nil},
		{start.Add(testSched[0]),
			[]string{"06"},
		},
		{start.Add(testSched[0] * 11),
			[]string{"08", "07", "06", "05"},
		},
		{start.Add(testSched[0] * 31),
			[]string{"01", "09", "02", "08", "03", "07", "04", "06"},
		},
	}

	for _, test := range tests {
		t.Run(test.now.String(), func(t *testing.T) {
			var calledFiles []string
			prune.Prune(fs, allFiles, test.now, testSched, func(fn string) {
				calledFiles = append(calledFiles, fn)
			})

			if !reflect.DeepEqual(test.expected, calledFiles) {
				t.Errorf("got %#v, want %#v", calledFiles, test.expected)
			}
		})
	}
}

func TestFilesInInterval(t *testing.T) {
	fs := fstest.MapFS{
		"1": &fstest.MapFile{ModTime: makeTime("2009-11-01 10:00:00")},
		"2": &fstest.MapFile{ModTime: makeTime("2009-11-01 10:01:00")},
		"3": &fstest.MapFile{ModTime: makeTime("2009-11-01 10:02:00")},
		"4": &fstest.MapFile{ModTime: makeTime("2009-11-01 10:03:00")},
		"5": &fstest.MapFile{ModTime: makeTime("2009-11-01 10:04:00")},
		"6": &fstest.MapFile{ModTime: makeTime("2009-11-01 11:05:00")},
		"7": &fstest.MapFile{ModTime: makeTime("2009-11-01 12:00:00")},
		"8": &fstest.MapFile{ModTime: makeTime("2009-11-01 13:00:00")},
	}

	intervalFilenames := prune.FilesInInterval(fs, must(fs.Glob("*")),
		makeTime("2009-11-01 10:50:00"),
		makeTime("2009-11-01 12:20:00"),
	)

	expectedFilenames := []string{"6", "7"}
	if !reflect.DeepEqual(expectedFilenames, intervalFilenames) {
		t.Errorf("got %v, want %v", intervalFilenames, expectedFilenames)
	}
}
