package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	var address string
	rootCmd := &cobra.Command{
		Use:     "weplanx",
		Version: "v0.0.1",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return cmd.Help()
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(address)
		},
	}
	rootCmd.PersistentFlags().StringVarP(&address, "server", "s", "0.0.0.0:9000", "http server address")
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
