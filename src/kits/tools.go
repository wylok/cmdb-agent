package kits

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"net"
	"os"
)

func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func Max(value []uint64) uint64 {
	var max uint64
	for _, val := range value {
		if val > max {
			max = val
		}
	}
	return max
}

func CheckFile(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}

func MapToJson(param map[string]interface{}) string {
	dataType, _ := json.Marshal(param)
	dataString := string(dataType)
	return dataString
}

func GetHostIp() []string {
	var HostIp []string
	addrs, err := net.InterfaceAddrs()
	if err == nil {
		for _, address := range addrs {
			// 检查ip地址判断是否回环地址
			if I, ok := address.(*net.IPNet); ok && !I.IP.IsLoopback() {
				if I.IP.To4() != nil {
					HostIp = append(HostIp, I.IP.String())
				}
			}
		}
	}
	return HostIp
}
