package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "digest",
	Short: "digest: Your recent google docs changes in your inbox",
	Long: `digest: Your recent google docs changes in your inbox
Almost all of command line flags can also be passed in as
environement variables with 'DIG_' prefix: e.g. auth-dir -> DIG_AUTH_DIR
	`,
}
