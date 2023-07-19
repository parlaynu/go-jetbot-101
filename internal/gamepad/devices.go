package gamepad

import (
	"os"
	"path/filepath"

	"github.com/holoplot/go-evdev"
)

type DeviceInfo struct {
	Driver string
	Name   string
}

func Devices() ([]*DeviceInfo, error) {
	var devices []*DeviceInfo

	ipath := "/dev/input"

	entries, err := os.ReadDir(ipath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		dpath := filepath.Join(ipath, entry.Name())

		d, err := evdev.Open(dpath)
		if err != nil {
			continue
		}

		name, _ := d.Name()
		d.Close()

		devices = append(devices, &DeviceInfo{
			Driver: dpath,
			Name:   name,
		})
	}

	return devices, nil
}
