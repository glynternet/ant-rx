package usbprompt

import (
	"fmt"

	"github.com/google/gousb"
	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
)

func UserSelectDevice(ctx *gousb.Context) (*DeviceDesc, error) {
	ds, err := getDeviceDescriptions(ctx)
	if err != nil {
		fmt.Printf("Error while listing devices: %s\n", err)
	}

	containsAnt := containsANTFunc()
	prompt := promptui.Select{
		Label: "Select Device",
		// items ordered with ANT-containing descriptions first
		Items: append(
			ds.matching(func(desc *DeviceDesc) bool {
				return containsAnt(desc.HumanReadable())
			}),
			ds.matching(func(desc *DeviceDesc) bool {
				return !containsAnt(desc.HumanReadable())
			})...,
		).humanReadables(),
	}

	_, result, err := prompt.Run()
	if err != nil {
		return nil, errors.Wrap(err, "running prompt")
	}

	fmt.Printf("You chose %q\n", result)
	chosen := ds.get(result)
	if chosen == nil {
		return nil, errors.Errorf("chosen device cannot be found: %s", result)
	}
	return chosen, nil
}
