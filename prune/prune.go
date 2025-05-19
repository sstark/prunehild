package prune

import (
	"fmt"
	"io/fs"
	"slices"
	"time"
	//"prunehild/schedules"
)

func modTime(fsys fs.FS, fn string) time.Time {
	st, _ := fs.Stat(fsys, fn)
	return st.ModTime()
}

func intervalDuration(sched []time.Duration, i int) time.Duration {
	if i == 0 {
		return 0
	}
	return sched[i] + intervalDuration(sched, i-1)
}

func Prune(fsys fs.FS, filenames []string, now time.Time, sched []time.Duration, callback func(fn string)) {
	const maxKeep = 2
	var filenamesDone []string

	// interval 0 does not need pruning, start with 1
	for i := len(sched) - 2; i > 0; i-- {
		interval := sched[i]
		intervalFilenames := FilesInInterval(fsys, filenames,
			now.Add(-intervalDuration(sched, i+1)),
			now.Add(-intervalDuration(sched, i)))

		pruneAgain := false
		if len(intervalFilenames) <= 2 {
			continue
		}

		// highest interval in sched
		if (i == len(sched)-2) && (len(intervalFilenames) > maxKeep) {
			callback(intervalFilenames[0])
			filenamesDone = append(filenamesDone, intervalFilenames[0])
			pruneAgain = true
		}

		// remaining intervals
		youngest := intervalFilenames[len(intervalFilenames)-1]
		secondYoungest := intervalFilenames[len(intervalFilenames)-2]
		dist := modTime(fsys, youngest).Sub(modTime(fsys, secondYoungest))
		if dist.Seconds() < interval.Seconds() && !slices.Contains(filenamesDone, youngest) {
			callback(youngest)
			filenamesDone = append(filenamesDone, youngest)
			pruneAgain = true
		}

		if pruneAgain {
			filenames = subtractSlices(filenames, filenamesDone)
			fmt.Printf("done: %#v\n", filenamesDone)
			Prune(fsys, filenames, now, sched, callback)
		}
	}
}

func subtractSlices(sl []string, subsl []string) (newsl []string) {
	for _, elem := range sl {
		if !slices.Contains(subsl, elem) {
			newsl = append(newsl, elem)
		}
	}
	return
}

func FilesInInterval(fsys fs.FS, filenames []string, from, to time.Time) (intervalFiles []string) {
	for _, filename := range filenames {
		f, _ := fs.Stat(fsys, filename)
		if f.ModTime().After(from) && f.ModTime().Before(to) {
			intervalFiles = append(intervalFiles, filename)
		}
	}
	return
}
