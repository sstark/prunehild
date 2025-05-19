package prune

import (
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

func Prune(fsys fs.FS, filenames []string, now time.Time, sched []time.Duration, keepMax int, callback func(fn string)) {
	// The highest real interval is actually the second to last, because the highest
	// does not need pruning by definition.
	highestInterval := len(sched) - 2
	// remaining will be mutated by the different branches of the recursion below.
	// It's a shared resource between recursion branches to make sure, "seen" items
	// are not considered more than once. Whenever a filename is pruned, the callback
	// should be called and the filename be removed from the remaining slice.
	remaining := slices.Clone(filenames)

	var pruneRemaining func()
	pruneRemaining = func() {
		for i := highestInterval; i > 0; i-- {
			interval := sched[i]
			intervalFilenames := filesInInterval(fsys, remaining,
				now.Add(-intervalDuration(sched, i+1)),
				now.Add(-intervalDuration(sched, i)))

			pruneAgain := false
			// We want to ensure we keep always at least 2 items per interval, so
			// se can measure the distance between the two.
			if len(intervalFilenames) <= 2 {
				continue
			}

			// pruneThis will prune a filename and set up the containing function for another
			// recursion level.
			pruneThis := func(fn string) {
				callback(fn)
				remaining = deleteValFromSlice(remaining, fn)
				pruneAgain = true
			}

			// In the highest interval we only need to check for the absolute number
			// of items and keep pruning until we are <= keepMax.
			if (i == highestInterval) && (len(intervalFilenames) > keepMax) {
				pruneThis(intervalFilenames[0])
			}

			// For each (except the highest) interval, determine the time distance between
			// the last and second-last item. If too small for this interval, prune the
			// latter.
			youngest, dist := getYoungestWithDist(fsys, intervalFilenames)
			if dist.Seconds() < interval.Seconds() {
				pruneThis(youngest)
			}

			if pruneAgain {
				pruneRemaining()
			}
		}
	}
	pruneRemaining()
}

// getYoungestWithDist determines the time difference between the last and second-last
// filename in the given list of filenames, in given file system.
func getYoungestWithDist(fsys fs.FS, filenames []string) (string, time.Duration) {
	youngest := filenames[len(filenames)-1]
	secondYoungest := filenames[len(filenames)-2]
	dist := modTime(fsys, youngest).Sub(modTime(fsys, secondYoungest))
	return youngest, dist
}

func deleteValFromSlice(sl []string, val string) []string {
	delIndex := slices.Index(sl, val)
	slNew := slices.Delete(sl, delIndex, delIndex+1)
	return slNew
}

func filesInInterval(fsys fs.FS, filenames []string, from, to time.Time) (intervalFiles []string) {
	for _, filename := range filenames {
		f, _ := fs.Stat(fsys, filename)
		if f.ModTime().After(from) && f.ModTime().Before(to) {
			intervalFiles = append(intervalFiles, filename)
		}
	}
	return
}
