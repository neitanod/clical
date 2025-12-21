package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "0.1.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show clical version",
	Long:  "Show current version of clical",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("clical version %s\n", version)
	},
}
