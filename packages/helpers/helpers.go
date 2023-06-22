// Package helpers 存放辅助方法
package helpers

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// Empty 类似于 PHP 的 empty() 函数
func Empty(val interface{}) bool {
	if val == nil {
		return true
	}
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.String, reflect.Array:
		return v.Len() == 0
	case reflect.Map, reflect.Slice:
		return v.Len() == 0 || v.IsNil()
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return reflect.DeepEqual(val, reflect.Zero(v.Type()).Interface())
}

// FirstElement 安全地获取 args[0]，避免 panic: runtime error: index out of range
func FirstElement(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return ""
}

// RandomNumber 生成长度为 length 随机数字字符串
func RandomNumber(length int) string {
	table := [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	b := make([]byte, length)
	n, err := io.ReadAtLeast(rand.Reader, b, length)
	if n != length {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

// RandomString 生成长度为 length 的随机字符串
func RandomString(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	letters := "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for i, v := range b {
		b[i] = letters[v%byte(len(letters))]
	}
	return string(b)
}

// MD5 生成字符串的 MD5 值
func MD5(str string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

// FormatBytes 格式化bytes
func FormatBytes(size float64) string {
	units := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}

	i := 0
	for ; size >= 1024 && i < len(units); i++ {
		size /= 1024
	}

	return fmt.Sprintf("%.2f %s", size, units[i])
}

// Cut 裁剪字符串
func Cut(begin, end, str string) string {
	b := utf8.RuneCountInString(str[:strings.Index(str, begin)]) + utf8.RuneCountInString(begin)
	e := utf8.RuneCountInString(str[:strings.Index(str, end)]) - b
	return string([]rune(str)[b : b+e])
}

// GetNetInfo 获取网络统计信息
func GetNetInfo() (uint64, uint64) {
	file, err := os.Open("/proc/net/dev")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	allRs := make(map[string][]string)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		if lineNumber < 3 {
			continue
		}

		line := strings.TrimSpace(scanner.Text())
		line = strings.Replace(line, ":", " ", -1)
		re := regexp.MustCompile("[ ]+")
		line = re.ReplaceAllString(line, " ")
		arr := strings.Split(line, " ")

		if len(arr) > 0 && arr[0] != "" {
			allRs[arr[0]+strconv.Itoa(lineNumber)] = []string{arr[0], arr[1], arr[9]}
		}
	}

	var keys []string
	for key := range allRs {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	tx := uint64(0)
	rx := uint64(0)

	for _, key := range keys {
		if strings.Contains(key, "lo") {
			continue
		}
		val := allRs[key]
		txValue, err := strconv.ParseUint(val[2], 10, 64)
		if err == nil {
			tx += txValue
		}
		rxValue, err := strconv.ParseUint(val[1], 10, 64)
		if err == nil {
			rx += rxValue
		}
	}

	return tx, rx
}

// MonitoringInfo 监控信息
type MonitoringInfo struct {
	CpuUse         float64 `json:"cpu_use"`
	Uptime         float64 `json:"uptime"`
	UptimePercent  float64 `json:"uptime_percent"`
	MemTotal       float64 `json:"mem_total"`
	MemUse         float64 `json:"mem_use"`
	MemUsePercent  float64 `json:"mem_use_percent"`
	SwapTotal      float64 `json:"swap_total"`
	SwapUse        float64 `json:"swap_use"`
	SwapUsePercent float64 `json:"swap_use_percent"`
	NetTx          uint64  `json:"net_tx"`
	NetRx          uint64  `json:"net_rx"`
}

// GetMonitoringInfo 获取监控数据
func GetMonitoringInfo() (MonitoringInfo, error) {
	var res MonitoringInfo

	// 网络流量
	netTx1, netRx1 := GetNetInfo()
	time.Sleep(time.Second)
	netTx2, netRx2 := GetNetInfo()
	res.NetTx = netTx2 - netTx1
	res.NetRx = netRx2 - netRx1

	// CPU 信息
	cpuInfoRaw, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return MonitoringInfo{}, err
	}
	physicalArr := make(map[string]struct{})
	var siblingsSum float64
	var re = regexp.MustCompile(`\d+\.\d+`)
	uptimeOutput, err := exec.Command("uptime").Output()
	if err != nil {
		return MonitoringInfo{}, err
	}
	uptimeValues := re.FindAllString(string(uptimeOutput), -1)
	uptime1, _ := strconv.ParseFloat(uptimeValues[0], 64)
	res.Uptime = uptime1
	processors := bytes.Split(cpuInfoRaw, []byte("\nprocessor"))
	rePhysical := regexp.MustCompile(`physical id\s*:\s(.*)`)
	reSiblings := regexp.MustCompile(`siblings\s*:\s(.*)`)
	for _, v := range processors {
		physical := rePhysical.FindSubmatch(v)
		siblings := reSiblings.FindSubmatch(v)
		if len(physical) > 1 {
			pid := string(physical[1])
			if _, found := physicalArr[pid]; !found {
				if len(siblings) > 1 {
					siblingsValue, _ := strconv.ParseFloat(string(siblings[1]), 64)
					siblingsSum += siblingsValue
				}
				physicalArr[pid] = struct{}{}
			}
		}
	}

	// CPU 使用率
	cpuUse := 0.1
	psOutput, err := exec.Command("ps", "aux").Output()
	if err != nil {
		return MonitoringInfo{}, err
	}
	cpuRaw := strings.Split(string(psOutput), "\n")
	pid := os.Getpid()
	for _, v := range cpuRaw {
		v = strings.TrimSpace(v)
		v = regexp.MustCompile(`\s+`).ReplaceAllString(v, " ")
		values := strings.Split(v, " ")
		if len(values) > 2 {
			p, _ := strconv.Atoi(values[1])
			if p == pid {
				continue
			}
			cpu, _ := strconv.ParseFloat(values[2], 64)
			cpuUse += cpu
		}
	}
	cpuUse = cpuUse / siblingsSum
	if cpuUse > 100 {
		cpuUse = 100
	}
	res.CpuUse = cpuUse

	// 内存使用率
	freeOutput, err := exec.Command("free", "-m").Output()
	if err != nil {
		return MonitoringInfo{}, err
	}
	memRaw := strings.Split(string(freeOutput), "\n")
	var memList, swapList string
	for _, v := range memRaw {
		if strings.Contains(v, "Mem") {
			memList = regexp.MustCompile(`\s+`).ReplaceAllString(v, " ")
		} else if strings.Contains(v, "Swap") {
			swapList = regexp.MustCompile(`\s+`).ReplaceAllString(v, " ")
		}
	}
	memArr := strings.Split(memList, " ")
	swapArr := strings.Split(swapList, " ")
	memTotal, _ := strconv.ParseFloat(memArr[1], 64)
	swapTotal, _ := strconv.ParseFloat(swapArr[1], 64)
	memUse, _ := strconv.ParseFloat(memArr[2], 64)
	swapUse, _ := strconv.ParseFloat(swapArr[2], 64)
	memUseP := (memUse / memTotal) * 100
	swapUseP := (swapUse / swapTotal) * 100
	uptime1P := uptime1 * 10
	if uptime1P > 100 {
		uptime1P = 100
	}

	res.MemTotal = memTotal
	res.MemUse = memUse
	res.MemUsePercent = memUseP
	res.SwapTotal = swapTotal
	res.SwapUse = swapUse
	res.SwapUsePercent = swapUseP
	res.UptimePercent = uptime1P

	return res, nil
}
