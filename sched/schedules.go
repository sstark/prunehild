package sched

import (
	"time"
)

const (
	Second = time.Second
	Minute = time.Minute
	Hour   = time.Hour
	Day    = Hour * 24
	Week   = Day * 7
	Month  = Week * 4
	Year   = Day * 365
	Long   = Year * 100
)

var BuiltinSchedules = map[string][]time.Duration{
	"ph0": []time.Duration{
		Second * 5,
		Second * 20,
		Second * 40,
		Second * 80,
		Long,
	},
	"ph1": []time.Duration{
		Day,
		Week,
		Month,
		Month * 4,
		Year,
		Long,
	},
	"ph2": []time.Duration{
		Day,
		Week,
		Month,
		Month * 4,
		Long,
	},
	"ph3": []time.Duration{
		Day,
		Week,
		Month,
		Long,
	},
}
