package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"prunehild/sched"
)

func init() {
	root_cmd.AddCommand(schedules_cmd)
}

var schedules_cmd = &cobra.Command{
	Use:   "schedules",
	Short: "List available schedules",
	Long: `
List all schedules that you can use for pruning.`,
	Run: func(cmd *cobra.Command, args []string) {
		for name, schedule := range sched.BuiltinSchedules {
			var descr string
			descr, ok := sched.BuiltinSchedulesDescription[name]
			if !ok {
				descr = "<no description>"
			}
			fmt.Printf("%s (%s) %v\n", name, descr, schedule)
		}
	},
}
