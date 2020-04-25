package usbprompt

import (
	"fmt"
	"strings"

	"github.com/google/gousb"
	"github.com/pkg/errors"
)

func DetectDevice(ctx *gousb.Context) (*DeviceDesc, error) {
	containsAnt := containsANTFunc()
	dds, err := getDeviceDescriptions(ctx)
	if err != nil {
		fmt.Printf("Error while listing devices: %s\n", err)
	}

	ants := dds.matching(func(desc *DeviceDesc) bool {
		return containsAnt(desc.HumanReadable())
	})

	if len(ants) == 0 {
		return nil, errors.New("Could not find any ANT USB candidates")
	}

	if len(ants) > 1 {
		return nil, fmt.Errorf("expected 1 ANT USB candidate but found %d: %s",
			len(ants),
			strings.Join(ants.humanReadables(), ", "))
	}

	return ants[0], nil
}
