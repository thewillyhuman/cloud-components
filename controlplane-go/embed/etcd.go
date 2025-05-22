package embed

import (
	"controlplane-go/config"
	"controlplane-go/internal/logging"
	"fmt"
	"go.etcd.io/etcd/server/v3/embed"
	"go.uber.org/zap"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

// StartEmbeddedEtcdWithConfig starts a configurable embedded etcd instance.
func StartEmbeddedEtcdWithConfig(nodeName, dataDir, initialCluster, clusterState, advertisePeer string) (*embed.Etcd, error) {
	log := logging.Logger

	log.Info("Starting embedded etcd with config",
		zap.String("nodeName", nodeName),
		zap.String("dataDir", dataDir),
		zap.String("initialCluster", initialCluster),
		zap.String("clusterState", clusterState),
		zap.String("advertisePeer", advertisePeer),
	)

	cfg := embed.NewConfig()
	cfg.Name = nodeName
	cfg.ClusterState = clusterState
	cfg.InitialCluster = initialCluster
	cfg.LogLevel = "error"
	cfg.LogOutputs = []string{"/dev/null"}

	fullDir := filepath.Join(dataDir, nodeName)
	if err := os.MkdirAll(fullDir, 0755); err != nil {
		log.Error("Failed to create data directory", zap.String("path", fullDir), zap.Error(err))
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}
	cfg.Dir = fullDir

	peerURL, err := url.Parse(advertisePeer)
	if err != nil {
		log.Error("Invalid peer URL", zap.String("url", advertisePeer), zap.Error(err))
		return nil, fmt.Errorf("invalid peer URL: %w", err)
	}
	cfg.ListenPeerUrls = []url.URL{*peerURL}
	cfg.AdvertisePeerUrls = []url.URL{*peerURL}

	clientURL, _ := url.Parse(fmt.Sprintf("http://127.0.0.1:2379"))
	cfg.ListenClientUrls = []url.URL{*clientURL}
	cfg.AdvertiseClientUrls = []url.URL{*clientURL}

	e, err := embed.StartEtcd(cfg)
	if err != nil {
		log.Error("Failed to start embedded etcd", zap.Error(err))
		return nil, err
	}

	select {
	case <-e.Server.ReadyNotify():
		log.Info("Embedded etcd server is ready", zap.String("nodeName", nodeName))
	case <-time.After(30 * time.Second):
		e.Server.Stop()
		log.Error("Etcd startup timeout", zap.String("nodeName", nodeName))
		return nil, fmt.Errorf("etcd startup timeout")
	}

	return e, nil
}

func StartEmbeddedEtcd() (*embed.Etcd, error) {
	log := logging.Logger

	log.Info("Starting embedded etcd")

	cfg := embed.NewConfig()
	cfg.Name = config.ControlPlaneName
	cfg.LogLevel = "error"
	cfg.LogOutputs = []string{"/dev/null"}

	fullDir := filepath.Join(config.EtcdDataDir)
	if err := os.MkdirAll(fullDir, 0755); err != nil {
		log.Error("Failed to create data directory", zap.String("path", fullDir), zap.Error(err))
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}
	cfg.Dir = fullDir

	lpurl, err := url.Parse(config.EtcdListenPeersUrl)
	if err != nil {
		log.Error("Invalid peer URL", zap.String("url", config.EtcdListenPeersUrl), zap.Error(err))
		return nil, fmt.Errorf("invalid peer URL: %w", err)
	}

	lcurl, err := url.Parse(config.EtcdListenClientsUrl)
	if err != nil {
		log.Error("Invalid client URL", zap.String("url", config.EtcdListenClientsUrl), zap.Error(err))
		return nil, fmt.Errorf("invalid client URL: %w", err)
	}

	cfg.ListenPeerUrls = []url.URL{*lpurl}
	cfg.AdvertisePeerUrls = []url.URL{*lpurl}
	cfg.ListenClientUrls = []url.URL{*lcurl}
	cfg.AdvertiseClientUrls = []url.URL{*lcurl}
	cfg.InitialCluster = fmt.Sprintf("%s=%s", config.ControlPlaneName, lpurl.String())
	cfg.ClusterState = "new"

	e, err := embed.StartEtcd(cfg)
	if err != nil {
		log.Error("Failed to start embedded etcd", zap.Error(err))
		return nil, err
	}

	select {
	case <-e.Server.ReadyNotify():
		log.Info("Etcd server is ready")
	case <-time.After(30 * time.Second):
		e.Server.Stop()
		log.Error("Etcd server took too long to start")
		return nil, fmt.Errorf("etcd server took too long to start")
	}

	return e, nil
}
