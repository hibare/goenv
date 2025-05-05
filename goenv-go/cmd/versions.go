package cmd

import (
	"fmt"

	"github.com/go-nv/goenv/internal/versions"
	"github.com/spf13/cobra"
)

var versionsCmd = &cobra.Command{
	Use:   "versions",
	Short: "List all installed Go versions",
	Long:  `List all Go versions that are currently installed.`,
	Run: func(cmd *cobra.Command, args []string) {
		vm, err := versions.NewVersionManager()
		if err != nil {
			fmt.Println(err)
			return
		}

		versions, err := vm.ListVersions()
		if err != nil {
			fmt.Println(err)
			return
		}

		// ToDo: Handle current version

		for _, version := range versions {
			fmt.Printf("  %s\n", version)
		}

	},
}

func init() {
	rootCmd.AddCommand(versionsCmd)
}
