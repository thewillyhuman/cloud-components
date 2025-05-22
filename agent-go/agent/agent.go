package agent

import "log"

const (
	NewStatus           = "New"
	CommissioningStatus = "Commissioning"
	ReadyStatus         = "Ready"
	TestingStatus       = "Testing"
	ControllerStatus    = "Controller"
	WorkerStatus        = "Worker"
)

var (
	status = NewStatus
)

func SetStatus(newStatus string) {
	log.Printf("Setting status to %s", newStatus)
	status = newStatus
}

func GetStatus() string {
	log.Printf("Fetching status: %s", status)
	return status
}
