package commissioning

import (
	"agent-go/agent"
	"fmt"
	//"github.com/shirou/gopsutil/v4/cpu"
	//"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
	"os/exec"
	"strings"
)

type Server struct {
	BMCInfo         BMCInfo
	NCores          int64
	MemoryInBytes   int64
	Disks           []Disk
	NICs            []NIC
	OperatingSystem string
}

type BMCInfo struct {
	IP string
}

type Disk struct {
	Device      string
	MountPoint  string
	SizeInBytes int64
}

type NIC struct {
	Device           string
	MacAddress       string
	IPAddress        string
	SpeedInMegabytes int64
}

func run(cmd string, args ...string) string {
	out, err := exec.Command(cmd, args...).Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(out))
}

func gatherBMCConfig() *BMCInfo {
	ipmi := run("ipmitool", "lan", "print")
	if strings.Contains(ipmi, "IP Address") {
		lines := strings.Split(ipmi, "\n")
		for _, line := range lines {
			if strings.Contains(line, "IP Address") {
				parts := strings.Split(line, ":")
				if len(parts) == 2 {
					ip := strings.TrimSpace(parts[1])
					return &BMCInfo{IP: ip}
				}
			}
		}
	}
	return nil
}

func GatherInventory() {
	agent.SetStatus(agent.CommissioningStatus)
	//BMCConfig := gatherBMCConfig()
	//nCores, _ := cpu.Counts(true)
	//memBytes, _ := mem.VirtualMemory()
	nics, _ := net.Interfaces()
	fmt.Printf("NICs: %s\n", nics)
}
