package exporter

import (
	"fmt"
//	"reflect"
)


func (Resource *ServerStatItem) Compare(number Number)(bool) {
	fmt.Println(Resource.DeviceLinkingCountInt, Resource.DiscReadFloat, Resource.DiscWriteFloat, Resource.MemoryUsageRate, Resource.DiscInfoList, Resource.CpuUsageRate, Resource.Load, Resource.IopsRead, Resource.IopsWrite)

	for _, diskInfo := range Resource.DiscInfoList {
		fmt.Println(diskInfo.DiscName, diskInfo.UsageRate)	
//		fmt.Println(reflect.TypeOf(diskInfo.DiscName))
//		fmt.Println(reflect.TypeOf(diskInfo.UsageRate))
		if (diskInfo.UsageRate) > number.DiskUsage {
			return true
		}
	}	
	
//    if (Resource.DeviceLinkingCountInt > number.DeviceLinkingCount) || (Resource.DiscReadFloat > number.DiscRead) || (Resource.DiscWriteFloat > number.DiscWrite) || (Resource.MemoryUsageRate > number.MemoryUsageRate) || (Resource.CpuUsageRate > number.CpuUsageRate) || (Resource.Load > number.Load) || (Resource.IopsRead > number.IopsRead) || (Resource.IopsWrite > numberIopsWrite) {
//		return true
//    } else {
//		return false
//	}


	switch {
	case Resource.DeviceLinkingCountInt > number.DeviceLinkingCount || Resource.DiscReadFloat > number.DiscRead || Resource.DiscWriteFloat > number.DiscWrite || Resource.MemoryUsageRate > number.MemoryUsageRate || Resource.CpuUsageRate > number.CpuUsageRate || Resource.Load > number.Load || Resource.IopsRead > number.IopsRead || Resource.IopsWrite > number.IopsWrite:
		return true

	default:
		return false
	}
}
