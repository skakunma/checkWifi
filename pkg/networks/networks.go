package networks

import (
	"bytes"
	"os/exec"
	"strconv"
	"strings"
)

type Network struct {
	BSSID           string
	Name            string
	Signal          int
	Channel         string
	HandshakeCap    string
	DecodedPassword string
}

type Networks = []Network

func GetNetworks(networkDevice string, limitNetworks int64) (Networks, error) {
	var resp Networks
	cmd := exec.Command("nmcli", "-t", "-f", "SSID,BSSID,SIGNAL", "dev", "wifi", "list", "ifname", networkDevice)
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return resp, err
	}

	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	for _, line := range lines {
		parts := strings.Split(line, ":")
		if len(parts) < 3 {
			continue
		}

		signal, err := strconv.Atoi(parts[2])
		if err != nil {
			continue
		}

		network := Network{
			Name:   parts[0],
			BSSID:  parts[1],
			Signal: signal,
		}
		resp = append(resp, network)
	}

	// Обрезаем до limit
	if limitNetworks > 0 && len(resp) > int(limitNetworks) {
		resp = resp[:limitNetworks]
	}

	return resp, nil
}

func FormatNetworks(networks Networks) []string {
	var resp []string
	for _, network := range networks {
		resp = append(resp, network.Name+" ("+strconv.FormatInt(int64(network.Signal), 10)+")")
	}
	return resp

}
