package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"nebulaLive/internal/app"
	"nebulaLive/internal/config"
	"os"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "nebula-live",
		Short: "Nebula Live CLI",
		Run: func(cmd *cobra.Command, args []string) {
			app.New().Run()
		},
	}
	config.InitConfig(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
