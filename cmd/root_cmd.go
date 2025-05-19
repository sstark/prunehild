package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var root_cmd = &cobra.Command{
	Use:   "prunehild",
	Short: "Prunehild is a file pruning utility",
	Long: `Use it to thin out directories full of regularly created files (for example backups),
with a defined schedule based on file modification time.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func Execute() {
	if err := root_cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
