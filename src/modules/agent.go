package modules

import (
	"bytes"
	_ "bytes"
	"cmdb-agent/config"
	"cmdb-agent/kits"
	"encoding/json"
	_ "github.com/CodyGuo/godaemon"
	"github.com/asmcos/requests"
	"github.com/mitchellh/go-ps"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/yumaojun03/dmidecode"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func HardWare() config.HardWareConf {
	//获取硬件信息
	conf := config.HardWareConf{}
	//获取CPU信息
	core, err := cpu.Counts(false)
	if err == nil {
		conf.CpuCore = core
	}
	info, err := cpu.Info()
	if err == nil {
		conf.CpuInfo = info[0].ModelName
	}
	//获取内存信息
	ms, err := mem.VirtualMemory()
	if err == nil {
		if ms != nil {
			conf.MemSize = ms.Total
		}
	}
	//获取磁盘信息
	ds := make(map[string]config.DiskInfo)
	dk, err := disk.Partitions(false)
	if err == nil {
		for _, k := range dk {
			di := config.DiskInfo{}
			di.FsType = k.Fstype
			di.MountPoint = k.Mountpoint
			d, err := disk.Usage(k.Mountpoint)
			if err == nil {
				di.Size = d.Total
			}
			ds[k.Device] = di
		}
		conf.Disk = ds
	}
	//获取网卡信息
	ns := make(map[string]config.NetInfo)
	nt, err := net.Interfaces()
	if err == nil {
		for _, k := range nt {
			ni := config.NetInfo{}
			ni.HardwareAddr = k.HardwareAddr
			for _, i := range k.Addrs {
				ni.Addrs = append(ni.Addrs, i.Addr)
			}
			ns[k.Name] = ni
		}
		conf.Net = ns
	}
	//获取bios信息
	dmi, err := dmidecode.New()
	if dmi != nil {
		infos, err := dmi.System()
		conf.Manufacturer = infos[0].Manufacturer
		conf.ProductName = infos[0].ProductName
		conf.SerialNumber = infos[0].SerialNumber
		if err == nil {
			conf.UUID = infos[0].UUID
		}
	}
	return conf
}

func HostSystem(hwConf config.HardWareConf) config.HostSystemConf {
	//获取操作系统信息
	conf := config.HostSystemConf{}
	h, err := host.Info()
	if err == nil {
		conf.Hostname = h.Hostname
		conf.Uptime = h.Uptime
		conf.OS = h.OS
		conf.Platform = h.Platform
		conf.PlatformVersion = h.PlatformVersion
		conf.KernelVersion = h.KernelVersion
		if strings.Contains(hwConf.Manufacturer, "VMware") {
			conf.HostID = kits.GetHostId(config.HostIdFile)
		} else {
			f, _ := os.Create(config.HostIdFile)
			_, _ = io.WriteString(f, h.HostID)
			conf.HostID = h.HostID
		}
	}
	return conf
}

func Collection() config.CollectionData {
	co := config.CollectionData{}
	co.HardWare = HardWare()
	co.System = HostSystem(co.HardWare)
	return co
}
func DownloadAgent(agent, AgentFile string) bool {
	// 文件重新下载
	res, err := http.Get(config.AgentFileUrl + agent)
	kits.Log("重新下载"+agent, "info", "CheckAgent")
	if err != nil {
		kits.Log(err.Error(), "error", "CheckAgent")
	} else {
		f, err := os.Create(AgentFile)
		if err != nil {
			kits.Log(err.Error(), "error", "CheckAgent")
		} else {
			_, err = io.Copy(f, res.Body)
			if err == nil {
				if kits.CheckFile(AgentFile) {
					_ = f.Chmod(0755)
					f.Close()
					return true
				}
			}
		}
	}
	return false
}

func CheckAgent(agent, AgentPath string, force bool) {
	AgentPid := config.PidPath + agent + ".pid"
	AgentFile := AgentPath + agent
	cmd := exec.Command(AgentFile, "-d")
	if kits.CheckFile(AgentPid) {
		f, err := ioutil.ReadFile(AgentPid)
		if err != nil {
			kits.Log(err.Error(), "error", "CheckAgent")
		} else {
			p, err := ps.FindProcess(int(kits.BytesToInt64(f)))
			if err != nil {
				kits.Log(err.Error(), "error", "CheckAgent")
			}
			if p == nil {
				if kits.CheckFile(AgentFile) {
					_ = cmd.Run()
				} else {
					if DownloadAgent(agent, AgentFile) {
						_ = cmd.Run()
					}
				}
			}
		}
	} else {
		if force {
			if kits.CheckFile(AgentFile) {
				_ = cmd.Run()
			} else {
				if DownloadAgent(agent, AgentFile) {
					_ = cmd.Run()
				}
			}
		}
	}
}

func Agent(cnf config.Config) {
	//req.Debug = 1
	time.Sleep(time.Duration(config.Interval) * time.Second)
	JsonData := config.ApiJson{}
	Co := config.AgentConf{}
	req := requests.Requests()
	Encrypt := kits.NewEncrypt([]byte(config.CryptKey), 16)
	co := Collection()
	data := Encrypt.EncryptString(co.CollectionDataString())
	jsonData, _ := json.Marshal(map[string]string{"agent_id": cnf.AgentId,
		"version": config.Version, "data": data})
	var stringBuilder bytes.Buffer
	stringBuilder.WriteString(kits.ExportUrl())
	stringBuilder.WriteString(config.UrlPath)
	resp, err := req.PostJson(stringBuilder.String(), string(jsonData))
	if err != nil {
		kits.Log(err.Error(), "error", "Agent")
	} else {
		err = resp.Json(&JsonData)
		if err != nil {
			kits.Log(err.Error(), "error", "Agent")
		} else {
			if JsonData.Success {
				v, _ := Encrypt.DecryptString(JsonData.Data.(string))
				err = json.Unmarshal(v, &Co)
				if err != nil {
					kits.Log(err.Error(), "error", "Agent")
				} else {
					// 恢复instance_id文件
					if Co.HostID != "" {
						f, _ := os.Create(config.HostIdFile)
						_, _ = io.WriteString(f, Co.HostID)
						f.Close()
					}
					// 版本升级
					Version, _ := strconv.Atoi(config.Version)
					NewVersion, _ := strconv.Atoi(Co.Version)
					if Co.Action == "upgrade" && Version < NewVersion {
						if kits.CheckFile("/sbin/cmdb_agent") {
							_ = os.Remove("/sbin/cmdb_agent")
						}
						os.Exit(0)
					}
				}
			}
		}
	}
}
