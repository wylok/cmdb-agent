package kits

import (
	"cmdb-agent/config"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"syscall"
)

func Log(Msg, MsgType string, FuncName string) {

	Prefix := map[string]string{"info": "[Info]", "error": "[Error]", "debug": "[Debug]"}
	_, err := os.Stat(config.LogFile)
	if err == nil {
		logFile, err := os.OpenFile(config.LogFile, syscall.O_RDWR|syscall.O_APPEND, 0666)
		if err == nil {
			defer logFile.Close()
			debugLog := log.New(logFile, FuncName+Prefix[MsgType], log.LstdFlags)
			debugLog.Println(Msg)
		}
	} else {
		logFile, err := os.Create(config.LogFile)
		if err == nil {
			defer logFile.Close()
			debugLog := log.New(logFile, FuncName+Prefix[MsgType], log.LstdFlags)
			debugLog.Println(Msg)
		}
	}
}

func GetConfig(c config.Config) (config.Config, error) {
	//通过读取文件获取
	var build strings.Builder
	build.WriteString(config.FilePath)
	build.WriteString(".system_manager/agent.yaml")
	yamlFile, err := ioutil.ReadFile(build.String())
	if err != nil {
		Log(err.Error(), "error", "GetConfig")
	}
	err = yaml.Unmarshal(yamlFile, &c)
	return c, err
}

func ExportUrl() string {
	Url := config.ProxyUrl
	for _, ip := range GetHostIp() {
		ip = strings.Join(strings.Split(ip, ".")[:3], ".") + ".1"
		tcpAddr := net.TCPAddr{IP: net.ParseIP(ip), Port: 12345}
		conn, err := net.DialTCP("tcp", nil, &tcpAddr)
		if err == nil {
			Url = "http://" + ip + ":12345"
			_ = conn.Close()
			break
		}
	}
	return Url
}

func GetHostId(File string) string {
	hostId := ""
	if CheckFile(File) {
		f, err := os.OpenFile(File, os.O_RDONLY, 0600)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			ContentByte, _ := ioutil.ReadAll(f)
			hostId = strings.TrimSpace(string(ContentByte))
		}
		if f != nil {
			defer f.Close()
		}
	}
	return hostId
}
