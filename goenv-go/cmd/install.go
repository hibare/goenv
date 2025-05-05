package cmd

import (
	"fmt"

	"github.com/go-nv/goenv/internal/installer"
	"github.com/go-nv/goenv/internal/versions"
	"github.com/spf13/cobra"
)

var (
	forceInstall bool
	listVersions bool
	skipExisting bool
)

func listAvailableVersions(cmd *cobra.Command, args []string) error {
	installer, err := installer.NewInstaller()
	if err != nil {
		return err
	}
	versions, err := installer.ListAvailableVersions()
	if err != nil {
		return err
	}

	fmt.Println("Available versions:")
	for _, version := range versions {
		fmt.Println(version)
	}
	return nil
}

func installVersion(cmd *cobra.Command, args []string) error {
	version := args[0]

	if !skipExisting {
		vm, err := versions.NewVersionManager()
		if err != nil {
			return err
		}

		if vm.IsVersionInstalled(version) {
			fmt.Printf("Go %s is already installed\n", version)
			return nil
		}
	}

	installer, err := installer.NewInstaller()
	if err != nil {
		return err
	}

	if err := installer.Install(version); err != nil {
		return err
	}

	fmt.Printf("Successfully installed Go %s\n", version)
	return nil
}

var installCmd = &cobra.Command{
	Use:   "install <version>",
	Short: "Install a specific version of Go",
	Long: `Install a specific version of Go.
The version should be in the format of X.Y.Z (e.g., 1.21.0).`,
	Run: func(cmd *cobra.Command, args []string) {
		if listVersions {
			if err := listAvailableVersions(cmd, args); err != nil {
				fmt.Println("Error listing available versions:", err)
			}
			return
		}

		if len(args) == 0 {
			fmt.Println("Error: version is required")
			return
		}

		if err := installVersion(cmd, args); err != nil {
			fmt.Println("Error installing version:", err)
			return
		}

	},
}

func init() {
	installCmd.Flags().BoolVarP(&forceInstall, "force", "f", false, "Force install even if the version is already installed")
	installCmd.Flags().BoolVarP(&listVersions, "list", "l", false, "List all available versions")
	installCmd.Flags().BoolVarP(&skipExisting, "skip-existing", "s", false, "Skip installation if the version is already installed")
	rootCmd.AddCommand(installCmd)
}
