package cmd

import (
	"fmt"

	"github.com/go-nv/goenv/internal/versions"
	"github.com/spf13/cobra"
)

var globalCmd = &cobra.Command{
	Use:   "global [version]",
	Short: "Set or show the global Go version",
	Long: `Set or show the global Go version.
If no version is specified, the current global version will be shown.
If a version is specified, it will be set as the global version.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		vm, err := versions.NewVersionManager()
		if err != nil {
			fmt.Println(err)
			return
		}

		if len(args) == 0 {
			version, err := vm.GetGlobalVersion()
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(version)
			return
		}

		err = vm.SetGlobalVersion(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Global version set to", args[0])
	},
}

func init() {
	rootCmd.AddCommand(globalCmd)
}
