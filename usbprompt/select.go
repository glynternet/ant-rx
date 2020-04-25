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
		fmt.Printf("Error while listing devices: %s", err)
	}

	prompt := promptui.Select{
		Label: "Select Device",
		Items: ds.humanReadables(),
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
