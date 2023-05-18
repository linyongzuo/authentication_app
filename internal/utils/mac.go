package utils

import (
	"net"
)

func GetMac() (mac string) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return mac
	}
	macs := make([]string, 0)
	for _, netInterface := range netInterfaces {
		macAddr := netInterface.HardwareAddr.String()
		if len(macAddr) == 0 {
			continue
		}

		macs = append(macs, macAddr)
	}
	return macs[0]
}
