package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
)

type IPInfo struct {
	Code int `json:"code"`
	Data IP  `json:"data`
}

type IP struct {
	Country   string `json:"country"`
	CountryId string `json:"country_id"`
	Area      string `json:"area"`
	AreaId    string `json:"area_id"`
	Region    string `json:"region"`
	RegionId  string `json:"region_id"`
	City      string `json:"city"`
	CityId    string `json:"city_id"`
	Isp       string `json:"isp"`
}

// 详细IP信息
func GetIpInfo() {
	external_ip := GetExternalIP() //获取公网IP
	external_ip = strings.Replace(external_ip, "\n", "", -1)
	fmt.Println("公网ip是: ", external_ip)
	fmt.Println("------Dividing Line------")
	ip := net.ParseIP(external_ip)
	if ip == nil {
		fmt.Println("您输入的不是有效的IP地址，请重新输入！")
	} else {
		result := GetIpInfoTabaoAPI(string(external_ip)) //调用淘宝API的json数据
		if result != nil {
			fmt.Println("国家：", result.Data.Country)
			fmt.Println("地区：", result.Data.Area)
			fmt.Println("城市：", result.Data.City)
			fmt.Println("运营商：", result.Data.Isp)
		}
	}
	fmt.Println("------Dividing Line------")

	GetIntranetIp() // 获取本机ip

	fmt.Println("------Dividing Line------")

	ip_int := inet_aton(net.ParseIP(external_ip))
	fmt.Println("Convert IPv4 address to decimal number(base 10) :", ip_int)

	ip_result := inet_ntoa(ip_int)
	fmt.Println("Convert decimal number(base 10) to IPv4 address:", ip_result)

	fmt.Println("------Dividing Line------")

	is_between := IpBetween(net.ParseIP("0.0.0.0"), net.ParseIP("255.255.255.255"), net.ParseIP(external_ip))
	fmt.Println("check result: ", is_between)

	fmt.Println("------Dividing Line------")
	is_public_ip := IsPublicIP(net.ParseIP(external_ip))
	fmt.Println("It is public ip: ", is_public_ip)

	is_public_ip = IsPublicIP(net.ParseIP("169.254.85.131"))
	fmt.Println("It is public ip: ", is_public_ip)

	fmt.Println("------Dividing Line------")

}

// 调用淘宝API来显示IP信息
func GetIpInfoTabaoAPI(ip string) *IPInfo {
	url := "http://ip.taobao.com/service/getIpInfo.php?ip="
	url += ip

	resp, err := http.Get(url)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	var result IPInfo
	if err := json.Unmarshal(out, &result); err != nil {
		return nil
	}

	return &result
}

// 假设本机只有一个ip和mac，多个可用[]string 实现
// 获取本机全部ip地址
func GetLocalIps() (ips []string) {
	interfaceAddr, err := net.InterfaceAddrs()
	if err != nil {
		//fmt.Println("get local ip failed")
		return ips
	}
	for _, address := range interfaceAddr {
		ipNet, isVailIpNet := address.(*net.IPNet)
		if isVailIpNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ips = append(ips, ipNet.IP.String())
			}
		}
	}
	//fmt.Println("ips = ", ips)
	return ips
}

// 获取本机ip
func GetIntranetIp() string {
	interfaceAddrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("get local ip failed")
		return ""
	}
	for _, address := range interfaceAddrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			// if IsPublicIP(ipnet.IP) == true { //只获取公网IP
			// 	return ipnet.IP.String()
			// }
			if ipnet.IP.To4() != nil {
				fmt.Println("ip:", ipnet.IP.String())
				//return ipnet.IP.String()
			}
		}
	}
	return ""
}

// 判断ip地址区间
func IpBetween(from net.IP, to net.IP, test net.IP) bool {
	if from == nil || to == nil || test == nil {
		fmt.Println("An ip input is nil") // or return an error!?
		return false
	}

	from16 := from.To16()
	to16 := to.To16()
	test16 := test.To16()
	if from16 == nil || to16 == nil || test16 == nil {
		fmt.Println("An ip did not convert to a 16 byte") // or return an error!?
		return false
	}

	if bytes.Compare(test16, from16) >= 0 && bytes.Compare(test16, to16) <= 0 {
		return true
	}
	return false
}

// 通过调用外网长期稳定的api接口,读取本机的公网IP地址
func GetExternalIP() string {
	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	return string(content)
}

// 通过调用内网dns服务器8.8.8.8:80端口,获取本机正在使用的内网IP地址
func GetPulicIP() string {
	conn, _ := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()
	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")
	return localAddr[0:idx]
}

// 判断是否是公网IP
func IsPublicIP(IP net.IP) bool {
	if IP.IsLoopback() || IP.IsLinkLocalMulticast() || IP.IsLinkLocalUnicast() {
		return false
	}
	if ip4 := IP.To4(); ip4 != nil {
		switch true {
		case ip4[0] == 10:
			return false
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return false
		case ip4[0] == 192 && ip4[1] == 168:
			return false
		default:
			return true
		}
	}
	return false
}

// ip地址string转int
func inet_aton(ipnr net.IP) int64 {
	bits := strings.Split(ipnr.String(), ".")
	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	var sum int64
	sum += int64(b0) << 24
	sum += int64(b1) << 16
	sum += int64(b2) << 8
	sum += int64(b3)
	return sum
}

// ip地址int转string
func inet_ntoa(ipnr int64) net.IP {
	var bytes [4]byte
	bytes[0] = byte(ipnr & 0xFF)
	bytes[1] = byte((ipnr >> 8) & 0xFF)
	bytes[2] = byte((ipnr >> 16) & 0xFF)
	bytes[3] = byte((ipnr >> 24) & 0xFF)
	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0])
}

// 获取本机的MAC地址 一般取第一个MAC地址作为唯一标识。
func GetLocalMac() (mac string) {
	// 获取本机的MAC地址
	interfaces, err := net.Interfaces()
	if err != nil {
		return ""
		//panic("Poor soul, here is what you got: " + err.Error())
	}
	//for _, inter := range interfaces {
	//fmt.Println(inter.Name)
	//}
	inter := interfaces[0]
	return inter.HardwareAddr.String() //获取本机MAC地址
}
