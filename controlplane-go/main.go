package main

import (
	"controlplane-go/cmd"
	"controlplane-go/internal/logging"
)

func main() {
	logging.Init()
	defer logging.Sync()

	// Since cmd.Execute() doesn’t return error, no need to capture anything
	cmd.Execute()
}
