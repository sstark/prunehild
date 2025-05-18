package schedules_test

import (
	"prunehild/schedules"
	"reflect"
	"testing"
	"time"
)

var testSchedules = schedules.ScheduleList{
	"testSched1": {schedules.Second, schedules.Second * 2, schedules.Second * 4, schedules.Second * 8, schedules.Second * 16, schedules.Long},
}

type offsetTestPair struct {
	i       int
	seconds int64
}

func TestScheduleOffset(t *testing.T) {
	var tests = []offsetTestPair{
		{0, 0},
		{1, 2},
		{2, 6},
		{3, 14},
		{4, 30},
		{5, 3153600030},
	}
	for _, pair := range tests {
		v := int64(testSchedules["testSched1"].Offset(pair.i).Seconds())
		if v != pair.seconds {
			t.Errorf("offset(%v) got %v, expected %v", pair.i, v, pair.seconds)
		}
	}
}

type goalTestPair struct {
	i    int
	goal int
}

func TestScheduleGoal(t *testing.T) {
	var tests = []goalTestPair{
		{0, 2},
		{1, 2},
		{2, 2},
		{3, 2},
		{4, 197100000},
	}
	for _, pair := range tests {
		v := testSchedules["testSched1"].Goal(pair.i)
		if v != pair.goal {
			t.Errorf("goal(%v) got %v, expected %v", pair.i, v, pair.goal)
		}
	}
}

func TestSchedulesAddFromFile(t *testing.T) {
	sched := schedules.ScheduleList{}
	err := sched.AddFromFile("../testdata/test-schedule.json")
	if err != nil {
		t.Fatal(err)
	}
	wanted := schedules.IntervalList{
		time.Hour * 24,
		time.Hour * 168,
		time.Hour * 672,
		time.Hour * 876000,
	}
	got := sched["test1"]
	if !reflect.DeepEqual(got, wanted) {
		t.Errorf("wanted %v, got %v", wanted, got)
	}
}
