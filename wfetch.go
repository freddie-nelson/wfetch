package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}

	info, err := getInfo(homeDir)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(formatInfo(info))
}

func getInfo(home string) ([]string, error) {
	out, err := exec.Command(fmt.Sprintf("systeminfo")).Output()
	if err != nil {
		return []string{}, err
	}

	splitOutput := strings.Split(string(out), "\r\n")
	info := make([]string, 0)

	for _, s := range splitOutput {
		if strings.Contains(s, "Hotfix(s):") {
			break
		} else {
			temp := strings.Split(s, ": ")
			if len(temp) == 1 {
				continue
			}

			info = append(info, strings.TrimSpace(temp[1]))
		}
	}

	return info, nil
}

func formatInfo(info []string) string {
	host := info[0]
	os := info[1]
	kernel := strings.Split(info[2], " N/A")[0]
	bootTime := bootTimeToUptime(info[10])
	memTotal := memToInt(info[23])
	memUsed := memTotal - memToInt(info[24])
	return fmt.Sprintf(
		`OS: %s
Host: %s
Kernel: %s
Uptime: %s
Memory: %v MB / %v MB`,
		os, host, kernel, bootTime, memUsed, memTotal)
}

func memToInt(memStr string) int {
	memInt, _ := strconv.Atoi(strings.Replace(memStr[:len(memStr)-3], ",", "", -1))
	return memInt
}

func bootTimeToUptime(bootTime string) string {
	layout := "01/02/2006, 15:04:05"
	date, _ := time.Parse(layout, bootTime)

	duration := time.Since(date)
	days := int(duration.Hours()) % 24
	hours := int(duration.Hours()) - (24 * days)
	minutes := int(duration.Minutes()) - int(duration.Hours())*60

	uptime := make([]string, 0, 3)
	if days != 0 {
		uptime = append(uptime, fmt.Sprintf("%v days", days))
	}
	if hours != 0 {
		uptime = append(uptime, fmt.Sprintf("%v hours", hours))
	}
	if minutes != 0 {
		uptime = append(uptime, fmt.Sprintf("%v mins", minutes))
	}

	return strings.Join(uptime, ", ")
}
