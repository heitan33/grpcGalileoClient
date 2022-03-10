package exporter

import (
	"fmt"
	"reflect"
)


func (Resource *ServerStatItem) Compare(number Number)(bool) {
	fmt.Println(Resource.DeviceLinkingCountInt, Resource.DiscReadFloat, Resource.DiscWriteFloat, Resource.MemoryUsageRate, Resource.DiscInfoList, Resource.CpuUsageRate, Resource.Load, Resource.IopsRead, Resource.IopsWrite)

//	var dealNotDisk, dealDisk bool = false
	for _, diskInfo := range Resource.DiscInfoList {
		fmt.Println(diskInfo.DiscName, diskInfo.UsageRate)	
		fmt.Println(reflect.TypeOf(diskInfo.DiscName))
		fmt.Println(reflect.TypeOf(diskInfo.UsageRate))
		if (diskInfo.UsageRate) > number.DiskUsage {
//			dealDisk = true
			return true
		}
	}	
	
    if (Resource.DeviceLinkingCountInt > number.DeviceLinkingCount) {
		return true
    } else {
		return false
	}
	
//	if (dealDisk) || (dealNotDisk) {
//		return true
//	} else {
//		return false
//	} 
}
