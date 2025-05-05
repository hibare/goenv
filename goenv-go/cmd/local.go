package cmd

import (
	"fmt"

	"github.com/go-nv/goenv/internal/versions"
	"github.com/spf13/cobra"
)

var localCmd = &cobra.Command{
	Use:   "local [version]",
	Short: "Set or show the local Go version",
	Long: `Set or show the local Go version.
If no version is specified, the current local version will be shown.
If a version is specified, it will be set as the local version.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		vm, err := versions.NewVersionManager()
		if err != nil {
			fmt.Println(err)
			return
		}

		// if no version is specified, show the local version
		if len(args) == 0 {
			version, err := vm.GetLocalVersion()
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(version)
			return
		}

		// if a version is specified, set the local version
		err = vm.SetLocalVersion(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Local version set to", args[0])
	},
}

func init() {
	rootCmd.AddCommand(localCmd)
}
