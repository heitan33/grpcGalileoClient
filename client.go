package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"os/exec"
	"reflect"
	"strconv"
	"strings"
	"time"

	"grpcClient/exporter"

	"gopkg.in/yaml.v2"
)

type myConf exporter.Conf

func (c *myConf) getConf() *myConf {
	yamlFile, err := ioutil.ReadFile("galileo.yaml")
	if err != nil {
		log.Println("yamlFile.Get err", err.Error())
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Println("Unmarshal: ", err.Error())
	}
	return c
}

type SystemInfo struct {
	Commandline string
}

func (s SystemInfo) getSysInfo() string {
	cmd := exec.Command("/bin/bash", "-c", s.Commandline)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("Error:can not obtain stdout pipe for command:%s\n", err)
	}

	if err := cmd.Start(); err != nil {
		log.Println("Error:The command is err,", err)
	}

	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Println("ReadAll Stdout:", err.Error())
	}

	if err := cmd.Wait(); err != nil {
		log.Println("wait:", err.Error())
	}

	result := string(bytes)
	return (result)
}

func commandStringHandle(command string, stat SystemInfo) float64 {
	stat.Commandline = command
	iopsReadStr := stat.getSysInfo()
	iopsReadStr = strings.Replace(iopsReadStr, "\n", "", -1)
	iopsReadStr = strings.Trim(iopsReadStr, " ")
	iopsReadFloat, err := strconv.ParseFloat(iopsReadStr, 64)
	if err != nil {
		log.Println("iopsReadFloat Stdout:", err.Error())
	}
	return (iopsReadFloat)
}

func serverState() string {
	var discReadFloat, discWriteFloat, memoryUsageRateFloat2, cpuUsageRate, load, bandwidthUpload, bandwidthDownload, iopsReadFloat, iopsWriteFloat float64
	var stat SystemInfo
	var memInt int64
	var resourceInfo string
	stat.Commandline = "sar -u 1 1|sed -n '4p'|awk '{print $NF}'"
	cpuStr := stat.getSysInfo()
	cpuStr = strings.Replace(cpuStr, "\n", "", -1)
	cpuFloat, _ := strconv.ParseFloat(cpuStr, 64)
	cpuUsageRate = 100 - cpuFloat
	cpuUsageRate, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", cpuUsageRate), 64)
	stat.Commandline = "free -m|grep 'available'"
	memStr := stat.getSysInfo()
	if len(memStr) == 0 {
		stat.Commandline = "free -m|grep -w 'cache'|awk '{print $NF}'"
		memStr = stat.getSysInfo()
		memStr = strings.Replace(memStr, "\n", "", -1)
		memInt, _ = strconv.ParseInt(memStr, 10, 64)
	} else {
		stat.Commandline = "free -m|sed -n '2p'|awk '{print $NF}'"
		memStr = stat.getSysInfo()
		memStr = strings.Replace(memStr, "\n", "", -1)
		memInt, _ = strconv.ParseInt(memStr, 10, 64)
	}

	stat.Commandline = "free -m|grep -i 'mem'|awk '{print $2}'"
	memTotalStr := stat.getSysInfo()
	memTotalStr = strings.Replace(memTotalStr, "\n", "", -1)
	memTotalInt, _ := strconv.ParseInt(memTotalStr, 10, 64)
	memoryUsageRateFloat := float64((float64(memTotalInt)-float64(memInt))/float64(memTotalInt)) * 100
	memoryUsageRateFloat2, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", memoryUsageRateFloat), 64)

	stat.Commandline = "uptime|awk -F, '{print $(NF-1)}'"
	cpuLoadStr := stat.getSysInfo()
	cpuLoadStr = strings.Replace(cpuLoadStr, "\n", "", -1)
	cpuLoadStr = strings.Trim(cpuLoadStr, " ")
	cpuLoadFloat, _ := strconv.ParseFloat(cpuLoadStr, 64)

	stat.Commandline = "grep 'model name' /proc/cpuinfo | wc -l"
	cpuCoreCountStr := stat.getSysInfo()
	cpuCoreCountStr = strings.Replace(cpuCoreCountStr, "\n", "", -1)
	cpuCoreCountStr = strings.Trim(cpuCoreCountStr, " ")
	cpuCoreCountFloat, _ := strconv.ParseFloat(cpuCoreCountStr, 64)
	loadFloat := cpuLoadFloat / cpuCoreCountFloat * 100
	load, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", loadFloat), 64)

	stat.Commandline = "ss -nt|grep -i 'est'|wc -l"
	deviceLinkingCountStr := stat.getSysInfo()
	deviceLinkingCountStr = strings.Replace(deviceLinkingCountStr, "\n", "", -1)
	deviceLinkingCountStr = strings.Trim(deviceLinkingCountStr, " ")
	deviceLinkingCountInt, err := strconv.ParseInt(deviceLinkingCountStr, 10, 64)
	if err != nil {
		log.Println("deviceLinkingCountInt Stdout:", err.Error())
	}

	command := "sar -b 1 1|sed -n '4p'|awk '{print $(NF-1)}'"
	discReadFloat = commandStringHandle(command, stat)
	discReadFloat = discReadFloat / 2

	command = "sar -b 1 1|sed -n '4p'|awk '{print $(NF)}'"
	discWriteFloat = commandStringHandle(command, stat)
	discWriteFloat = discWriteFloat / 2

	command = "sar -b 1 1|sed -n '$p'|awk '{print $3}'"
	iopsReadFloat = commandStringHandle(command, stat)

	command = "sar -b 1 1|sed -n '$p'|awk '{print $4}'"
	iopsWriteFloat = commandStringHandle(command, stat)

	stat.Commandline = "df -P -h|awk '{print $1,$2,$(NF-1)}'|sed '1d'|grep '/.*'"
	diskStr := strings.TrimSpace(stat.getSysInfo())
	var diskStatSlice []string
	var diskName string
	var discInfo exporter.DiskInfo
	piontDisc := &discInfo
	var diskTotalFloat float64
	var discInfoList []exporter.DiskInfo
	diskStr = strings.Replace(diskStr, "\n", ",", -1)
	diskStatSlice = strings.Split(diskStr, ",")
	for _, diskStat_dev := range diskStatSlice {
		diskName = strings.Split(diskStat_dev, " ")[0]
		if strings.Contains(diskName, ".") {
			diskName = strings.Replace(diskName, ".", "_", -1)
		}
		diskTotal := strings.Split(diskStat_dev, " ")[1]
		if strings.Contains(diskTotal, "T") == true {
			diskTotal = diskTotal[0 : len(diskTotal)-1]
			diskTotalFloat, _ = strconv.ParseFloat(diskTotal, 32)
			diskTotalFloat = diskTotalFloat * 1024
		} else if strings.Contains(diskTotal, "M") == true {
			diskTotal = diskTotal[0 : len(diskTotal)-1]
			diskTotalFloat, _ = strconv.ParseFloat(diskTotal, 32)
			diskTotalFloat = diskTotalFloat / 1024
			diskTotalFloat, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", diskTotalFloat), 64)
		} else {
			diskTotal = diskTotal[0 : len(diskTotal)-1]
			diskTotalFloat, _ = strconv.ParseFloat(diskTotal, 32)
		}
		diskContent := strings.Split(diskStat_dev, " ")[2]
		diskContent = diskContent[0 : len(diskContent)-1]
		diskContentFloat, _ := strconv.ParseFloat(diskContent, 32)
		piontDisc.DiscName = diskName
		piontDisc.Total = diskTotalFloat
		piontDisc.UsageRate = diskContentFloat
		discInfoList = append(discInfoList, discInfo)
	}
	var c myConf
	conf := c.getConf()
	var hostName string
	hostName = conf.HostName

	resourceInfo = serverStatJson(hostName, discInfoList, deviceLinkingCountInt, discReadFloat, discWriteFloat, memoryUsageRateFloat2, cpuUsageRate, load, bandwidthUpload, bandwidthDownload, iopsReadFloat, iopsWriteFloat)
	log.Println(time.Now().String())
	return resourceInfo
}

type ServerStatItem exporter.ServerStatItem

func serverStatJson(hostName string, discInfoList []exporter.DiskInfo, deviceLinkingCountInt int64, numbers ...float64) string {
	serverStatItem := ServerStatItem{HostName: hostName, DiscInfoList: discInfoList, DeviceLinkingCountInt: deviceLinkingCountInt, DiscReadFloat: numbers[0], DiscWriteFloat: numbers[1], MemoryUsageRate: numbers[2], CpuUsageRate: numbers[3], Load: numbers[4], BandwidthUpload: numbers[5], BandwidthDownload: numbers[6], IopsRead: numbers[7], IopsWrite: numbers[8]}
	postItemJson, err := json.Marshal(serverStatItem)
	postItemStr := string(postItemJson)
	if err != nil {
		log.Println(err)
	}
	return postItemStr
}

//var resourceInfo ServerStatItem

var address string
var warning, midValue bool = false, false

func main() {
	var c myConf
	var count int8
	var Number exporter.Number
	var pioResourceInfo *exporter.ServerStatItem

	pioResourceInfo = new(exporter.ServerStatItem)
	pioNumber := &Number
	conf := c.getConf()

	pioNumber.CpuUsageRate = conf.CpuUsageRate
	pioNumber.MemoryUsageRate = conf.MemoryUsageRate
	pioNumber.DiskUsage = conf.DiskUsage
	pioNumber.Load = conf.Load
	pioNumber.DiscRead = conf.DiscRead
	pioNumber.DiscWrite = conf.DiscWrite
	pioNumber.IopsRead = conf.IopsRead
	pioNumber.IopsWrite = conf.IopsWrite
	pioNumber.DeviceLinkingCount = conf.DeviceLinkingCount

	address = conf.Host

	warning = pioResourceInfo.Tag

	for {
		serverState := serverState()
		fmt.Println(serverState)
		json.Unmarshal([]byte(serverState), &exporter.ResourceInfo)
		fmt.Println(exporter.ResourceInfo)
		fmt.Println(reflect.TypeOf(exporter.ResourceInfo))
		pioResourceInfo = &exporter.ResourceInfo
		//		warning = pioResourceInfo.Tag
		//		diskInfo := pioResourceInfo.DiscInfoList
		//		linkCount := pioResourceInfo.DeviceLinkingCountInt
		//		diskPerRead := pioResourceInfo.DiscReadFloat
		//		diskPerWrite := pioResourceInfo.DiscWriteFloat
		//		fmt.Println(diskInfo, linkCount, diskPerRead, diskPerWrite)
		fmt.Println(pioResourceInfo.Tag, pioResourceInfo.DiscInfoList, pioResourceInfo.DeviceLinkingCountInt, pioResourceInfo.DiscReadFloat, pioResourceInfo.DiscWriteFloat)

		//		resourceInfo.DeviceLinkingCountInt = DeviceLinkingCountInt
		//		var number int64 = 20
		warning := exporter.ResourceInfo.Compare(Number)
		fmt.Println("----------")
		fmt.Println(warning)
		fmt.Println("----------")
		if warning != midValue {
			count++
			if count == 3 {
				midValue = warning
				pioResourceInfo.Tag = midValue
				strResourceInfo, err := json.Marshal(exporter.ResourceInfo)
				if err != nil {
					fmt.Println(err)
				}
				exporter.Report(midValue, string(strResourceInfo), address)
			}
		} else {
			count = 0
		}
		time.Sleep(time.Duration(90) * time.Second)
	}
}
