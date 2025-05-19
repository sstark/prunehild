package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/fs"
	"os"
	"prunehild/prune"
	"prunehild/sched"
	"time"
)

var (
	// var Version = "1.0.0"
	fGlob     string
	fSchedule string
	fKeepmax  int
)

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
List all filenames that should be pruned according to the given schedule.
Only files matched by the glob pattern are considered.`,
	Run: func(cmd *cobra.Command, args []string) {
		//log.Printf("prunhild version %s\n", Version)
		schedule := sched.BuiltinSchedules[fSchedule]
		fsys := os.DirFS(".")
		fileNames, _ := fs.Glob(fsys, fGlob)
		prune.Prune(fsys, fileNames, time.Now(), schedule, fKeepmax, func(fn string) {
			fmt.Println(fn)
		})
	},
}
