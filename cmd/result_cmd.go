package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/fs"
	"os"
	"prunehild/prune"
	"prunehild/sched"
	"slices"
	"time"
)

func init() {
	root_cmd.AddCommand(result_cmd)
	result_cmd.Flags().StringVarP(&fGlob, "glob", "g", "*", "Filename pattern")
	result_cmd.Flags().StringVarP(&fSchedule, "schedule", "s", "ph1", "Schedule")
	result_cmd.Flags().IntVarP(&fKeepmax, "keepmax", "k", 2, "Max files to keep in highest interval")
}

var result_cmd = &cobra.Command{
	Use:   "result",
	Short: "List files not pruned",
	Long: `
List all filenames that would be unaffected by the given schedule.
This is like running 'ls' with the pruning candidates subtracted.`,
	Run: func(cmd *cobra.Command, args []string) {
		schedule := sched.BuiltinSchedules[fSchedule]
		fsys := os.DirFS(".")
		fileNames, _ := fs.Glob(fsys, fGlob)
		prune.Prune(fsys, fileNames, time.Now(), schedule, fKeepmax, func(fn string) {
			idx := slices.Index(fileNames, fn)
			fileNames = slices.Delete(fileNames, idx, idx+1)
		})
		for _, fn := range fileNames {
			fmt.Printf("%s\n", fn)
		}
	},
}
