package config

import "fmt"

const (
	// Etcd defaults
	DefaultEtcdDataDir              = "/var/lib/controlplane/etcd/"
	DefaultEtcdListenClientsAddress = "127.0.0.1"
	DefaultEtcdListenClientsPort    = 2379
	DefaultEtcdListenPeersAddress   = "0.0.0.0"
	DefaultEtcdListenPeersPort      = 2380
	DefaultUIListenAddress          = "0.0.0.0"
	DefaultUIListenPort             = 8080

	// Etcd control plane prefixes
	EtcdControlPlanePrefix      = "/controlplane"
	EtcdControlPlaneNodesPrefix = "/controlplane/%s/nodes"
)

var (
	EtcdDataDir          = DefaultEtcdDataDir
	EtcdListenClientsUrl = fmt.Sprintf("http://%s:%d", DefaultEtcdListenClientsAddress, DefaultEtcdListenClientsPort)
	EtcdListenPeersUrl   = fmt.Sprintf("http://%s:%d", DefaultEtcdListenPeersAddress, DefaultEtcdListenPeersPort)

	ControlPlaneName   = ""
	ControlPlaneRegion = ""
	AdvertiseAddress   = ""
	JoinPeerAddress    = ""

	UIListenUrl = fmt.Sprintf("%s:%d", DefaultUIListenAddress, DefaultUIListenPort)
)
