package cmd

import (
	"controlplane-go/internal/logging"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "controlplane",
	Short: "Control plane for distributed infrastructure",
}

func Execute() {
	log := logging.Logger

	log.Info("Welcome to Cloud Components Control Plane. Starting up...")

	cobra.OnInitialize()
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(joinCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Error("CLI execution failed", zap.Error(err))
		os.Exit(1)
	}
}
