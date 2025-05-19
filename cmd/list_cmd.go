package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/fs"
	"os"
	"prunehild/prune"
	"prunehild/schedules"
	"time"
)

var (
	// var Version = "1.0.0"
	fGlob     string
	fSchedule string
	fKeepmax  int
)

var builtinSchedules = map[string][]time.Duration{
	"ph0": []time.Duration{
		schedules.Second * 5,
		schedules.Second * 20,
		schedules.Second * 40,
		schedules.Second * 80,
		schedules.Long,
	},
	"ph1": []time.Duration{
		schedules.Day,
		schedules.Week,
		schedules.Month,
		schedules.Month * 4,
		schedules.Year,
		schedules.Long,
	},
	"ph2": []time.Duration{
		schedules.Day,
		schedules.Week,
		schedules.Month,
		schedules.Month * 4,
		schedules.Long,
	},
	"ph3": []time.Duration{
		schedules.Day,
		schedules.Week,
		schedules.Month,
		schedules.Long,
	},
}

func init() {
	root_cmd.AddCommand(list_cmd)
	list_cmd.Flags().StringVarP(&fGlob, "glob", "g", "*", "Filename pattern")
	list_cmd.Flags().StringVarP(&fSchedule, "schedule", "s", "ph1", "Schedule")
	list_cmd.Flags().IntVarP(&fKeepmax, "keepmax", "k", 2, "Max files to keep in highest interval")
}

var list_cmd = &cobra.Command{
	Use:   "list",
	Short: "List prune candidates",
	Long: `
FIXME`,
	Run: func(cmd *cobra.Command, args []string) {
		//log.Printf("prunhild version %s\n", Version)
		schedule := builtinSchedules[fSchedule]
		fsys := os.DirFS(".")
		fileNames, _ := fs.Glob(fsys, fGlob)
		prune.Prune(fsys, fileNames, time.Now(), schedule, fKeepmax, func(fn string) {
			fmt.Println(fn)
		})
	},
}
