package control

import (
	"context"
	"controlplane-go/config"
	"controlplane-go/embed"
	"controlplane-go/internal/logging"
	"controlplane-go/types"
	"controlplane-go/util"
	"controlplane-go/web"
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"os"
	"time"
)

func InitControlPlane() {
	log := logging.Logger

	log.Info(fmt.Sprintf("Using control plane config: name=%s, region=%s, etcd_data_dir=%s, advertise_address=%s, etcd_curl=%s, etcd_purl=%s",
		config.ControlPlaneName,
		config.ControlPlaneRegion,
		config.EtcdDataDir,
		config.AdvertiseAddress,
		config.EtcdListenClientsUrl,
		config.EtcdListenPeersUrl,
	))

	// Start the embed version of Etcd.
	etcdServer, err := embed.StartEmbeddedEtcd()
	if err != nil {
		log.Fatal("Failed to start embedded etcd", zap.Error(err))
	}
	<-etcdServer.Server.ReadyNotify()

	// Once the embed etcd is ready create a client to it.
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{config.EtcdListenClientsUrl},
		DialTimeout: 3 * time.Second,
	})
	if err != nil {
		log.Fatal("Failed to connect to etcd", zap.Error(err))
	}
	defer etcdClient.Close()

	// Store control plane metadata
	cpMeta := types.ControlPlaneMetadata{
		Name:   config.ControlPlaneName,
		Region: config.ControlPlaneRegion,
	}

	cpMetaData, err := json.Marshal(cpMeta)
	if err != nil {
		log.Fatal("Failed to serialize control plane metadata", zap.Error(err))
	}

	_, err = etcdClient.Put(context.Background(),
		fmt.Sprintf("/controlplane/%s/metadata", config.ControlPlaneName),
		string(cpMetaData),
	)
	if err != nil {
		log.Fatal("Failed to store control plane metadata", zap.Error(err))
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal("Failed to get hostname", zap.Error(err))
	}

	meta := util.DetectNodeMetadata()

	node := types.NodeInfo{
		Hostname:  hostname,
		IP:        config.AdvertiseAddress,
		Provider:  meta.Provider,
		Location:  meta.Location,
		NodeID:    meta.NodeID,
		OSName:    meta.OSName,
		OSVersion: meta.OSVersion,
		Labels: map[string]string{
			"role": "control-plane",
		},
	}

	nodeData, err := json.Marshal(node)
	if err != nil {
		log.Fatal("Failed to serialize node info", zap.Error(err))
	}

	_, err = etcdClient.Put(context.Background(),
		fmt.Sprintf("/controlplane/%s/nodes/%s", config.ControlPlaneName, hostname),
		string(nodeData),
	)
	if err != nil {
		log.Fatal("Failed to store node info", zap.Error(err))
	}

	log.Info("Control plane initialized and node registered")

	web.StartUI(config.ControlPlaneName, hostname, config.AdvertiseAddress, "http://127.0.0.1:2379")

	<-etcdServer.Server.StopNotify()
	log.Info("Etcd server stopped. Exiting control plane node.")
}
