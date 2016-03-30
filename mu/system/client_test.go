package system

import (
	"testing"
)

func TestGetUptime(t *testing.T) {
	output, err := GetUptime()
	if err != nil {
		t.Error(err)
	}
	t.Log(output)
}

func TestGetLoad(t *testing.T) {
	output, err := GetLoad()
	if err != nil {
		t.Error(err)
	}
	t.Log(output)
}
