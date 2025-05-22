package web

import (
	"context"
	"controlplane-go/config"
	"controlplane-go/internal/logging"
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"html/template"
	"net/http"
	"time"
)

type NodeInfo struct {
	Hostname  string            `json:"hostname"`
	IP        string            `json:"ip"`
	Provider  string            `json:"provider"`
	Location  string            `json:"location"`
	NodeID    string            `json:"node_id"`
	OSName    string            `json:"os_name"`
	OSVersion string            `json:"os_version"`
	Labels    map[string]string `json:"labels"`
}

type PageData struct {
	ControlPlane string
	Node         string
	IP           string
	Peers        []NodeInfo
}

func StartUI(cpName, nodeName, ip string, etcdEndpoint string) {
	log := logging.Logger

	tpl := template.Must(template.New("ui").Parse(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>Control Plane UI</title>
			<style>
				body { font-family: sans-serif; padding: 1rem; }
				table { border-collapse: collapse; width: 100%; margin-top: 1rem; }
				th, td { border: 1px solid #ddd; padding: 0.5rem; text-align: left; }
				th { background-color: #f2f2f2; }
			</style>
		</head>
		<body>
			<h2>Control Plane: {{.ControlPlane}}</h2>
			<p><strong>This Node:</strong> {{.Node}} ({{.IP}})</p>
			<h3>All Nodes:</h3>
			<table>
				<tr>
					<th>Hostname</th>
					<th>IP</th>
					<th>Provider</th>
					<th>Location</th>
					<th>Node ID</th>
					<th>OS</th>
					<th>Labels</th>
				</tr>
				{{range .Peers}}
				<tr>
					<td>{{.Hostname}}</td>
					<td>{{.IP}}</td>
					<td>{{.Provider}}</td>
					<td>{{.Location}}</td>
					<td>{{.NodeID}}</td>
					<td>{{.OSName}} {{.OSVersion}}</td>
					<td>{{.Labels}}</td>
				</tr>
				{{else}}
				<tr><td colspan="6">No nodes found</td></tr>
				{{end}}
			</table>
		</body>
		</html>
	`))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Info("Serving UI request")

		cli, err := clientv3.New(clientv3.Config{
			Endpoints:   []string{etcdEndpoint},
			DialTimeout: 3 * time.Second,
		})
		if err != nil {
			http.Error(w, "etcd connection failed", 500)
			log.Error("Failed to connect to etcd", zap.Error(err))
			return
		}
		defer cli.Close()

		keyPrefix := fmt.Sprintf("/controlplane/%s/nodes/", cpName)
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		resp, err := cli.Get(ctx, keyPrefix, clientv3.WithPrefix())
		if err != nil {
			http.Error(w, "failed to query etcd", 500)
			log.Error("Failed to get peers from etcd", zap.Error(err))
			return
		}

		var peers []NodeInfo
		for _, kv := range resp.Kvs {
			var node NodeInfo
			if err := json.Unmarshal(kv.Value, &node); err == nil {
				peers = append(peers, node)
			} else {
				log.Warn("Failed to decode node info", zap.String("key", string(kv.Key)), zap.Error(err))
			}
		}

		data := PageData{
			ControlPlane: cpName,
			Node:         nodeName,
			IP:           ip,
			Peers:        peers,
		}

		if err := tpl.Execute(w, data); err != nil {
			log.Error("Failed to render template", zap.Error(err))
		}
	})

	go func() {
		// Update the listen interface for the UI with the address used to advertise the control plane.
		config.UIListenUrl = fmt.Sprintf("%s:%d", config.AdvertiseAddress, config.DefaultUIListenPort)
		log.Info(fmt.Sprintf("Starting node web UI: http://%s", config.UIListenUrl))
		if err := http.ListenAndServe(config.UIListenUrl, nil); err != nil {
			log.Fatal("Web UI failed", zap.Error(err))
		}
	}()
}
