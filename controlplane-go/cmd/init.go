package cmd

import (
	"controlplane-go/config"
	"controlplane-go/control"
	"controlplane-go/internal/logging"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new control plane",
	Run: func(cmd *cobra.Command, args []string) {
		log := logging.Logger

		log.Info(
			fmt.Sprintf("Initializing control plane [%s]: region=%s, etcdDataDir=%s, advertise_address=%s.",
				config.ControlPlaneName,
				config.ControlPlaneRegion,
				config.EtcdDataDir,
				config.AdvertiseAddress,
			))

		if config.ControlPlaneName == "" || config.ControlPlaneRegion == "" || config.AdvertiseAddress == "" {
			log.Error("Flags --name, --region, and --advertise-ip are required")
			os.Exit(1)
		}

		control.InitControlPlane()
	},
}

func init() {
	initCmd.Flags().StringVar(&config.ControlPlaneName, "name", "", "Name of the control plane")
	initCmd.Flags().StringVar(&config.ControlPlaneRegion, "region", "", "Region of the control plane")
	initCmd.Flags().StringVar(&config.EtcdDataDir, "data-dir", "", "Directory to store control plane data")
	initCmd.Flags().StringVar(&config.AdvertiseAddress, "advertise-ip", "", "IP address of the node to advertise the control plane")
}
