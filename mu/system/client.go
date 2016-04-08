package system

import (
	"os/exec"
	"strings"
)

type Client struct {
}

func GetLoad() (string, error) {
	output, err := exec.Command("cat", "/proc/loadavg").Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func GetUptime() (string, error) {
	output, err := exec.Command("cat", "/proc/uptime").Output()
	if err != nil {
		return "", err
	}
	loadAry := strings.Split(string(output), " ")
	return loadAry[0], nil
}
