package cmd

import (
	"controlplane-go/config"
	"controlplane-go/control"
	"controlplane-go/internal/logging"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var joinCmd = &cobra.Command{
	Use:   "join",
	Short: "Join an existing control plane cluster",
	Run: func(cmd *cobra.Command, args []string) {
		log := logging.Logger

		log.Info(
			fmt.Sprintf("Joining control plane [%s]: etcdDataDir=%s, advertise_address=%s",
				config.JoinPeerAddress,
				config.EtcdDataDir,
				config.AdvertiseAddress,
			))

		if config.JoinPeerAddress == "" || config.AdvertiseAddress == "" {
			log.Error("Flags --peer-ip and --advertise-ip are required")
			os.Exit(1)
		}

		control.JoinControlPlane()
	},
}

func init() {
	joinCmd.Flags().StringVar(&config.JoinPeerAddress, "peer-ip", "", "Peer address of existing cluster")
	joinCmd.Flags().StringVar(&config.EtcdDataDir, "data-dir", "/etc/controlplane/data", "Directory to store control plane data")
	joinCmd.Flags().StringVar(&config.AdvertiseAddress, "advertise-ip", "", "IP address of the node to advertise the control plane")
}
