
package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "godict",
}

func Execute() {
	rootCmd.Execute()
}