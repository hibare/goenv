package cmd

import (
	"fmt"

	"github.com/go-nv/goenv/internal/installer"
	"github.com/go-nv/goenv/internal/versions"
	"github.com/spf13/cobra"
)

var (
	forceUninstall bool
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall <version>",
	Short: "Uninstall a specific version of Go",
	Long: `Uninstall a specific version of Go.
The version should be in the format of X.Y.Z (e.g., 1.21.0).`,
	Args:    cobra.ExactArgs(1),
	Example: "goenv uninstall <version>",
	Run: func(cmd *cobra.Command, args []string) {
		version := args[0]

		vm, err := versions.NewVersionManager()
		if err != nil {
			fmt.Println("Error creating version manager:", err)
			return
		}

		if !vm.IsVersionInstalled(version) {
			fmt.Printf("version %s is not installed\n", version)
			return
		}

		if !forceUninstall {
			fmt.Printf("Are you sure you want to uninstall Go %s? (y/n): ", version)
			var confirm string
			fmt.Scanln(&confirm)
			if confirm != "y" {
				fmt.Println("uninstall cancelled")
				return
			}
		}

		installer, err := installer.NewInstaller()
		if err != nil {
			fmt.Println("Error creating installer:", err)
			return
		}

		if err := installer.Uninstall(version); err != nil {
			fmt.Printf("failed to uninstall Go %s: %s\n", version, err)
			return
		}

		fmt.Printf("Successfully uninstalled Go %s\n", version)
	},
}

func init() {
	uninstallCmd.Flags().BoolVarP(&forceUninstall, "force", "f", false, "Force uninstall without confirmation")
	rootCmd.AddCommand(uninstallCmd)
}
