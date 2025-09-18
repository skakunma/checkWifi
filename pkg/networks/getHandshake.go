package networks

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func (n *Network) LoadHandshake(networkDevice string) error {
	startMonCmd := exec.Command("bash", "-c", fmt.Sprintf("sudo airmon-ng start %s", networkDevice))
	if output, err := startMonCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("ошибка запуска режима мониторинга: %v\n%s", err, string(output))
	}

	monitorDevice := networkDevice + "mon"

	outputDir := "handshakes"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("не удалось создать директорию %s: %w", outputDir, err)
	}

	fileBase := fmt.Sprintf("%s/%s_%d", outputDir, sanitizeFileName(n.Name), time.Now().Unix())
	capPath := fileBase + "-01.cap"
	n.HandshakeCap = capPath

	captureCmd := fmt.Sprintf("sudo timeout 30s airodump-ng --bssid %s --channel %s --write %s %s",
		n.BSSID, n.Channel, fileBase, monitorDevice)

	fmt.Println("Захват handshake...")
	cmd := exec.Command("bash", "-c", captureCmd)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("ошибка при захвате handshake: %v\n%s", err, string(output))
	}

	if _, err := os.Stat(capPath); os.IsNotExist(err) {
		return fmt.Errorf("файл handshake не найден: %s", capPath)
	}

	checkCmd := fmt.Sprintf("aircrack-ng %s -a2 -b %s", capPath, n.BSSID)
	out, err := exec.Command("bash", "-c", checkCmd).CombinedOutput()
	if err != nil {
		return fmt.Errorf("ошибка при проверке handshake: %v\n%s", err, string(out))
	}

	if !strings.Contains(strings.ToLower(string(out)), "handshake") {
		return fmt.Errorf("handshake не обнаружен в файле: %s", capPath)
	}

	fmt.Println("✅ Handshake успешно перехвачен и сохранён в:", capPath)
	return nil
}

func sanitizeFileName(name string) string {
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "/", "_")
	name = strings.ReplaceAll(name, "\\", "_")
	return name
}
