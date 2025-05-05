package cmd

import (
	"fmt"
	"os"

	"github.com/go-nv/goenv/internal/utils"
	"github.com/go-nv/goenv/internal/version"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "goenv",
	Short: "Go version manager",
	Long: `A simple and powerful Go version manager.
It allows you to easily switch between multiple versions of Go.`,
	Version: version.CurrentVersion,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	if err := utils.InitDirs(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
