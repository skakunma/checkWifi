package networkDevices

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type NetworkInterface struct {
	Name  string
	Speed int
}

type NetworkDevices = []NetworkInterface

func FormatNetworkDevices(networkDevices NetworkDevices) ([]string, error) {
	var resp = make([]string, 0)
	for _, device := range networkDevices {
		resp = append(resp, device.Name+fmt.Sprintf(" SPEED: (%d)", device.Speed))
	}
	return resp, nil
}

func GetTopNetworkDevicesLinux() NetworkDevices {
	devices := NetworkDevices{}

	interfaces, err := ioutil.ReadDir("/sys/class/net/")
	if err != nil {
		fmt.Println("Ошибка при чтении сетевых интерфейсов:", err)
		return devices
	}

	for _, iface := range interfaces {
		ifaceName := iface.Name()
		speedPath := filepath.Join("/sys/class/net/", ifaceName, "speed")

		speedBytes, err := ioutil.ReadFile(speedPath)
		if err != nil {
			continue
		}

		speedStr := strings.TrimSpace(string(speedBytes))
		speed, err := strconv.Atoi(speedStr)
		if err != nil {
			continue
		}

		devices = append(devices, NetworkInterface{
			Name:  ifaceName,
			Speed: speed,
		})
	}

	sort.Slice(devices, func(i, j int) bool {
		return devices[i].Speed > devices[j].Speed
	})

	return devices
}
