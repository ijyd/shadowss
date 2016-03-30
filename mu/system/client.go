package system

import (
	"os/exec"
)

type Client struct {
}

func GetLoad() (string, error) {
	uptimeCommand := `cat /proc/loadavg | awk '{ print $1" "$2" "$3 }'`
	output, err := exec.Command(uptimeCommand).Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func GetUptime() (string, error) {
	uptimeCommand := `cat /proc/uptime | awk '{ print $1 }'`
	output, err := exec.Command(uptimeCommand).Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
