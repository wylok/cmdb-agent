package config

const (
	Version    = "2021042306"
	Interval   = 60
	UrlPath    = "/xxx"
	CryptKey   = "4098879a2529ca11b8675505ahf88a2d"
	FilePath   = "/xxx"
	HostIdFile = FilePath + "instance-id"
	PidPath    = FilePath + "/"
	PidFile    = PidPath + "cmdb_agent.pid"
	LogFile    = PidPath + "cmdb_agent.log"
	ProxyUrl   = "http://xxx:12345"
)

type Config struct {
	AgentId string `yaml:"agent_id"`
}

type ApiJson struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type DiskInfo struct {
	MountPoint string `json:"mount_point"`
	FsType     string `json:"fs_type"`
	Size       uint64 `json:"size"`
}

type NetInfo struct {
	HardwareAddr string   `json:"hardwareaddr"`
	Addrs        []string `json:"addrs"`
}

type HardWareConf struct {
	UUID         string              `json:"uuid"`
	SerialNumber string              `json:"serial_number"`
	Manufacturer string              `json:"manufacturer"`
	ProductName  string              `json:"product_name"`
	CpuCore      int                 `json:"cpu_core"`
	CpuInfo      string              `json:"cpu_info"`
	MemSize      uint64              `json:"mem_size"`
	Disk         map[string]DiskInfo `json:"disk"`
	Net          map[string]NetInfo  `json:"net"`
}

type HostSystemConf struct {
	HostID          string `json:"host_id"`
	Hostname        string `json:"hostname"`
	OS              string `json:"os"`
	Platform        string `json:"platform"`
	PlatformVersion string `json:"platformVersion"`
	KernelVersion   string `json:"kernelVersion"`
	Uptime          uint64 `json:"uptime"`
}

type CollectionData struct {
	HardWare HardWareConf   `json:"hardware"`
	System   HostSystemConf `json:"system"`
}

type AgentConf struct {
	Action  string `json:"action"`
	HostID  string `json:"host_id"`
	Version string `json:"version"`
}
