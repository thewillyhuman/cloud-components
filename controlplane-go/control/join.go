package control

import (
	"context"
	"controlplane-go/config"
	"controlplane-go/embed"
	"controlplane-go/internal/logging"
	"controlplane-go/types"
	"controlplane-go/util"
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"os"
	"strings"
	"time"
)

func JoinControlPlane() {
	log := logging.Logger

	log.Info(fmt.Sprintf("Using control plane config: peer=%s, etcd_data_dir=%s, advertise_address=%s, etcd_curl=%s, etcd_purl=%s",
		config.JoinPeerAddress,
		config.EtcdDataDir,
		config.AdvertiseAddress,
		config.EtcdListenClientsUrl,
		config.EtcdListenPeersUrl,
	))

	// Step 1: Connect to peer etcd
	log.Info(fmt.Sprintf("Connecting to peer=%s", config.JoinPeerAddress))
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{config.JoinPeerAddress},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal("Failed to connect to peer etcd", zap.Error(err))
	} else {
		log.Info(fmt.Sprintf("Connected to peer=%s", config.JoinPeerAddress))
	}
	defer cli.Close()

	// Step 2: Discover control plane name
	log.Info(fmt.Sprintf("Downloading control plane data from peer=%s", config.JoinPeerAddress))
	metaPrefix := "/controlplane/"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := cli.Get(ctx, metaPrefix, clientv3.WithPrefix(), clientv3.WithKeysOnly())
	if err != nil {
		log.Fatal("Failed to discover control plane keys", zap.Error(err))
	}

	var cpName string
	for _, kv := range resp.Kvs {
		key := string(kv.Key)
		if strings.HasSuffix(key, "/metadata") {
			parts := strings.Split(key, "/")
			cpName = parts[2]
			break
		}
	}

	if cpName == "" {
		log.Fatal("Could not identify control plane name from peer")
	}

	log.Info("Discovered control plane",
		zap.String("controlPlane", cpName),
	)

	// Step 3: Get existing peers list
	peerListKey := fmt.Sprintf("/controlplane/%s/peers", cpName)
	peerResp, err := cli.Get(context.Background(), peerListKey)
	if err != nil {
		log.Fatal("Failed to get peer list", zap.Error(err))
	}

	var peers []string
	if len(peerResp.Kvs) > 0 {
		_ = json.Unmarshal(peerResp.Kvs[0].Value, &peers)
	}

	// Step 4: Get local IP and hostname
	ip, err := util.GetOutboundIP()
	if err != nil {
		log.Fatal("Failed to get local IP", zap.Error(err))
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal("Failed to get hostname", zap.Error(err))
	}

	thisPeerURL := fmt.Sprintf("http://0.0.0.0:2380")

	log.Info("Local node info",
		zap.String("hostname", hostname),
		zap.String("ip", ip),
	)

	// Step 5: Add this node to the etcd cluster
	addResp, err := cli.MemberAdd(context.Background(), []string{thisPeerURL})
	if err != nil {
		log.Fatal("Failed to add member to etcd", zap.Error(err))
	}

	// Step 6: Build initial cluster string
	clusterStr := ""
	for _, m := range addResp.Members {
		for _, u := range m.PeerURLs {
			clusterStr += fmt.Sprintf("%s=%s,", m.Name, u)
		}
	}
	clusterStr = strings.TrimRight(clusterStr, ",")

	log.Info("Built initial cluster string", zap.String("cluster", clusterStr))

	// Step 7: Start embedded etcd node
	etcd, err := embed.StartEmbeddedEtcdWithConfig(
		hostname,
		config.EtcdDataDir,
		clusterStr,
		"existing",
		thisPeerURL,
	)
	if err != nil {
		log.Fatal("Failed to start embedded etcd node", zap.Error(err))
	}
	<-etcd.Server.ReadyNotify()

	log.Info("Etcd node started and joined the cluster")

	// Step 8: Register this node locally
	localCli, _ := clientv3.New(clientv3.Config{
		Endpoints: []string{config.EtcdListenClientsUrl},
	})
	defer localCli.Close()

	node := types.NodeInfo{
		Hostname: hostname,
		IP:       ip,
	}

	nodeKey := fmt.Sprintf("/controlplane/%s/nodes/%s", cpName, hostname)
	nodeData, _ := json.Marshal(node)

	_, err = localCli.Put(context.Background(), nodeKey, string(nodeData))
	if err != nil {
		log.Fatal("Failed to register node in etcd", zap.Error(err))
	}

	log.Info("Registered node in etcd",
		zap.String("key", nodeKey),
		zap.String("ip", ip),
	)

	// Step 9: Update peer list
	peers = append(peers, ip)
	peerData, _ := json.Marshal(peers)

	_, err = localCli.Put(context.Background(), peerListKey, string(peerData))
	if err != nil {
		log.Fatal("Failed to update peer list", zap.Error(err))
	}

	log.Info("Updated peer list",
		zap.String("controlPlane", cpName),
		zap.Strings("peers", peers),
	)

	log.Info("Node successfully joined control plane",
		zap.String("hostname", hostname),
		zap.String("controlPlane", cpName),
	)

	<-etcd.Server.StopNotify()
	log.Info("Etcd server stopped. Node exiting.")
}
