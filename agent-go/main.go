package main

import (
	"agent-go/commissioning"
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

const configPath = "/tmp/node-agent.json"

type BMCInfo struct {
	IP      string `json:"ipmi_ip,omitempty"`
	Present bool   `json:"present"`
}

type Inventory struct {
	CPU   int      `json:"cpu"`
	RAM   string   `json:"ram"`
	Disks []string `json:"disks"`
	Brand string   `json:"brand"`
	BMC   BMCInfo  `json:"bmc"`
}

type Config struct {
	Mode      string    `json:"mode"` // not_enrolled, controller, worker
	UUID      string    `json:"uuid"`
	Inventory Inventory `json:"inventory"`
}

func run(cmd string, args ...string) string {
	out, err := exec.Command(cmd, args...).Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(out))
}

func detectBMC() BMCInfo {
	ipmi := run("ipmitool", "lan", "print")
	if strings.Contains(ipmi, "IP Address") {
		lines := strings.Split(ipmi, "\n")
		for _, line := range lines {
			if strings.Contains(line, "IP Address") {
				parts := strings.Split(line, ":")
				if len(parts) == 2 {
					ip := strings.TrimSpace(parts[1])
					return BMCInfo{IP: ip, Present: true}
				}
			}
		}
	}
	return BMCInfo{Present: false}
}

func gatherInventory() Inventory {
	osType := runtime.GOOS
	var cores int
	var ram, brand string
	var disks []string
	var bmc BMCInfo

	switch osType {
	case "linux":
		cores, _ = cpu.Counts(true)
		ram = run("free", "-h")
		brand = run("dmidecode", "-s", "system-manufacturer")
		disksRaw := run("lsblk", "-d", "-o", "NAME,SIZE")
		disks = strings.Split(disksRaw, "\n")
		bmc = detectBMC()
	case "darwin":
		cores, _ = cpu.Counts(true)
		memBytes, _ := mem.VirtualMemory()
		ram = formatMacMemory(strconv.FormatUint(memBytes.Total, 10))
		brand = run("system_profiler", "SPHardwareDataType")
		disksRaw := run("diskutil", "list")
		disks = strings.Split(disksRaw, "\n")
		bmc = BMCInfo{Present: false}
	default:
		cores = -1
		ram = "unknown"
		disks = []string{"unknown"}
		brand = "unknown"
		bmc = BMCInfo{Present: false}
	}

	return Inventory{
		CPU:   cores,
		RAM:   ram,
		Disks: disks,
		Brand: brand,
		BMC:   bmc,
	}
}

func formatMacMemory(bytesStr string) string {
	bytesStr = strings.TrimSpace(bytesStr)
	if bytesStr == "" {
		return "unknown"
	}

	var bytesVal int64
	fmt.Sscanf(bytesStr, "%d", &bytesVal)
	gb := float64(bytesVal) / (1024 * 1024 * 1024)
	return fmt.Sprintf("%.1f GB", gb)
}

func getUUID() string {
	if runtime.GOOS == "linux" {
		return run("dmidecode", "-s", "system-uuid")
	} else if runtime.GOOS == "darwin" {
		out := run("ioreg", "-rd1", "-c", "IOPlatformExpertDevice")
		for _, line := range strings.Split(out, "\n") {
			if strings.Contains(line, "IOPlatformUUID") {
				parts := strings.Split(line, "\"")
				if len(parts) > 3 {
					return parts[3]
				}
			}
		}
	}
	return "unknown"
}

func loadOrCreateConfig() Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		inv := gatherInventory()
		uuid := run("dmidecode", "-s", "system-uuid")
		config := Config{
			Mode:      "not_enrolled",
			UUID:      uuid,
			Inventory: inv,
		}
		saveConfig(config)
		return config
	} else {
		data, _ := ioutil.ReadFile(configPath)
		var cfg Config
		_ = json.Unmarshal(data, &cfg)
		return cfg
	}
}

func saveConfig(cfg Config) {
	data, _ := json.MarshalIndent(cfg, "", "  ")
	_ = ioutil.WriteFile(configPath, data, 0644)
}

func main() {
	//cfg := loadOrCreateConfig()
	//fmt.Printf("Node Agent running in mode: %s\n", cfg.Mode)
	//fmt.Printf("Inventory:\n%+v\n", cfg.Inventory)

	//disks, _ := disk.Partitions(false)
	//fmt.Printf("Disks: %+v\n", disks)
	commissioning.GatherInventory()
	commissioning.RunBenchmarks()

	// Future: Add enroll logic, controller registration, etc.
}
