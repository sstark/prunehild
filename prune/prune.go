package prune

import (
	"io/fs"
	//"prunehild/schedules"
)

func Prune(fs fs.StatFS, filenames []string, callback func(fn string)) {
	//intervals := schedules.Schedules["shortterm"]
	// interval 0 does not need pruning, start with 1
	//for i := len(intervals) - 2; i > 0; i-- {
	// find filenames in fs
	// return if found less then minimum files (e. g. 2) as safety
	// iv := find files that fall into interval i
	// pruneAgain := false
	// if len(iv) <= 2 { continue }
	callback("3")
	callback("4")
}
