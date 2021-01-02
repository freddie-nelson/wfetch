package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/go-ps"
)

// SysInfo is a store for system info
type SysInfo struct {
	user           string
	host           string
	os             string
	kernel         string
	bootTime       string
	shell          string
	de             string
	wm             string
	terminal       string
	cpu            string
	gpu            string
	memUsed        int
	memTotal       int
	terminalHeight int
}

// FormatInfo formats the given SysInfo into a multiline string
func FormatInfo(info SysInfo) string {
	userAtHost := fmt.Sprintf("%s%s@%s%s", info.user, ResetAnsii, AccentAnsii, info.host)
	underline := ResetAnsii + strings.Repeat("-", len(userAtHost)-len(AccentAnsii)-len(ResetAnsii))

	s := fmt.Sprintf(
		`%s
%s
OS: %s
Host: %s
Kernel: %s
Uptime: %s
DE: %s
WM: %s
Shell: %s
Terminal: %s
CPU: %s
GPU: %s
Memory: %v MB / %v MB`,
		userAtHost, underline, info.os, info.host, info.kernel, info.bootTime, info.de, info.wm, info.shell, info.terminal, info.cpu, info.gpu, info.memUsed, info.memTotal)

	return strings.Replace(s, ": ", ResetAnsii+": ", -1)
}

// GetInfo gets system info for the current users environment
func GetInfo() SysInfo {
	info := SysInfo{}

	// get computersystem info
	computerSystemInfo := getValuesFromList(exec.Command("wmic", "computersystem", "get", "name,username", "/value").Output())
	info.host = computerSystemInfo[0]
	info.user = strings.Split(computerSystemInfo[1], "\\")[1]

	// get os info
	osInfo := getValuesFromList(exec.Command("wmic", "os", "get", "caption,freephysicalmemory,lastbootuptime,totalvisiblememorysize,version", "/value").Output())
	info.os = osInfo[0]
	info.kernel = osInfo[4]
	info.bootTime = bootTimeToUptime(osInfo[2])
	info.memTotal = memToInt(osInfo[3])
	info.memUsed = info.memTotal - memToInt(osInfo[1])

	// hardware info
	info.cpu = getValuesFromList(exec.Command("wmic", "cpu", "get", "name", "/value").Output())[0]
	if strings.Contains(info.cpu, "Intel") {
		info.cpu = info.cpu[18:]
	}

	info.gpu = getValuesFromList(exec.Command("wmic", "path", "win32_VideoController", "get", "caption", "/value").Output())[0]

	// windows specific values
	info.de = "Aero"
	info.wm = "Explorer"

	// terminal
	// terminalLines, _ := exec.Command("mode", "con", "/status", "|", "findstr", "Lines:").Output()
	// info.terminalHeight, _ = strconv.Atoi(strings.Split(string(terminalLines), ":")[1])

	pProcess, _ := ps.FindProcess(os.Getppid())
	ppProcess, _ := ps.FindProcess(pProcess.PPid())
	if ppProcess.Executable() == "explorer.exe" {
		info.terminal = pProcess.Executable()[:len(pProcess.Executable())-4]
		info.shell = getShell(info.terminal, info.kernel)
	} else {
		info.terminal = ppProcess.Executable()[:len(ppProcess.Executable())-4]
		info.shell = getShell(pProcess.Executable()[:len(pProcess.Executable())-4], info.kernel)
	}

	return info
}

func getShell(term string, kern string) string {
	switch term {
	case "cmd":
		return fmt.Sprintf("%s %s", term, kern)
	case "pwsh":
		fallthrough
	case "powershell":
		version := getPowershellVersion()

		return fmt.Sprintf("%s %s", term, version)
	}

	return ""
}

func getPowershellVersion() string {
	major, _ := exec.Command("powershell", "(Get-Variable PSVersionTable -ValueOnly).PSVersion.Major").Output()
	minor, _ := exec.Command("powershell", "(Get-Variable PSVersionTable -ValueOnly).PSVersion.Minor").Output()
	patch, _ := exec.Command("powershell", "(Get-Variable PSVersionTable -ValueOnly).PSVersion.Patch").Output()

	versionSlice := make([]string, 0, 3)
	if len(major) > 0 {
		versionSlice = append(versionSlice, string(major))
	}
	if len(minor) > 0 {
		versionSlice = append(versionSlice, string(minor))
	}
	if len(patch) > 0 {
		versionSlice = append(versionSlice, string(patch))
	}

	var version string
	for _, s := range versionSlice {
		version += "." + strings.Trim(s, "\r\n")
	}

	return version[1:]
}

func getValuesFromList(out []byte, _ error) []string {
	list := strings.Split(string(out), "\r")

	values := make([]string, 0)
	for _, s := range list {
		if strings.Contains(s, "=") {
			values = append(values, strings.Split(s, "=")[1])
		}
	}

	return values
}

func bootTimeToUptime(bootTime string) string {
	layout := "20060201150405"
	bootTimeWMI := strings.Split(bootTime, ".")[0]
	bootTimeParsed, _ := time.Parse(layout, bootTimeWMI)

	duration := time.Since(bootTimeParsed)
	days := int(duration.Hours() / 24)
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

func memToInt(memStr string) int {
	memInt, _ := strconv.Atoi(memStr)
	return memInt / 1024
}
