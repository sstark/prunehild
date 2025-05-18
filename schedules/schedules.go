package schedules

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"time"
)

const (
	Second               = time.Second
	Minute               = time.Minute
	Hour                 = time.Hour
	Day                  = Hour * 24
	Week                 = Day * 7
	Month                = Week * 4
	Year                 = Day * 365
	Long                 = Year * 100
	DefaultSchedFileName = "/etc/prunehild.schedules"
)

type IntervalList []time.Duration

// Offset returns how Long ago the given interval started
func (il IntervalList) Offset(i int) time.Duration {
	if i == 0 {
		return 0
	}
	return il[i] + il.Offset(i-1)
}

// Goal returns how many snapshots are the Goal in the given interval
func (il IntervalList) Goal(i int) int {
	if i > len(il)-2 {
		panic("this should not happen: highest interval is innumerable!")
	}
	return int(il[i+1] / il[i])
}

type ScheduleList map[string]IntervalList

type jsonInterval []map[string]time.Duration

func (schl *ScheduleList) String() string {
	var a []string
	for sch := range *schl {
		a = append(a, sch)
	}
	return strings.Join(a, ",")
}

// Schedules is a list of available snapshot Schedules. Defines how often
// snapshots are made or purged. The span of an interval is always the snapshot
// distance of the next interval.
var Schedules = ScheduleList{
	"longterm":  {Hour * 6, Day, Week, Month, Long},
	"shortterm": {Minute * 10, Hour * 2, Day, Week, Month, Long},
}

// AddFromFile adds an external JSON file to the list of available scheds
func (schl ScheduleList) AddFromFile(file string) error {
	// If we are using the default file name, and it doesn't exist, no problem, just return
	if _, err := os.Stat(file); os.IsNotExist(err) && file == DefaultSchedFileName {
		return nil
	}
	schedFile, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("Error opening schedule file: %v", err)
	}
	var readData map[string]jsonInterval
	err = json.Unmarshal(schedFile, &readData)
	if err != nil {
		return fmt.Errorf("Error parsing schedule file: %v", err)
	}
	for k, v := range readData {
		schl[k] = v.intervalList()
	}
	return nil
}

// list prints the stored Schedules in the list
func (schl ScheduleList) list() {
	var sKeys []string
	for k := range schl {
		sKeys = append(sKeys, k)
	}
	sort.Strings(sKeys)
	for _, name := range sKeys {
		fmt.Printf("%s: %s\n", name, schl[name])
	}
}

// Transform a JSON formatted intervalList like this:
// [
//   { "day" : 1, "Hour" : 12 },
//   { "week" : 2 },
//   { "month" : 1, "week" : 2}
//   { "long" : 1}
// ]
// and it makes it equivalent to
// { 1*day + 12*Hour, 2*week, 1*month + 2*week, long }

func (json jsonInterval) intervalList() IntervalList {
	il := make(IntervalList, len(json))
	for i, interval := range json {
		var duration time.Duration
	Loop:
		for k, v := range interval {
			switch k {
			case "Second":
				duration += v * Second
			case "Minute":
				duration += v * Minute
			case "Hour":
				duration += v * Hour
			case "Day":
				duration += v * Day
			case "Week":
				duration += v * Week
			case "Month":
				duration += v * Month
			case "Year":
				duration += v * Year
			case "Long":
				duration = Long
				break Loop
			}
		}
		il[i] = duration
	}
	return il
}
