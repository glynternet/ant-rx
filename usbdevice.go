package main

import (
	"fmt"
	"regexp"

	"github.com/google/gousb"
	"github.com/google/gousb/usbid"
	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
)

type deviceDesc gousb.DeviceDesc

func (dd deviceDesc) humanReadable() string {
	return fmt.Sprintf("%s (%s) (%s)",
		usbid.Vendors[dd.Vendor].Product[dd.Product],
		usbid.Vendors[dd.Vendor],
		usbid.Classes[dd.Class])
}

type deviceDescs []*deviceDesc

func (dds deviceDescs) humanReadables() []string {
	pattern := regexp.MustCompile(`.*[aA][nN][tT].*`)
	var containingAnt []string
	var items []string
	for _, d := range dds {
		desc := d.humanReadable()
		if pattern.MatchString(desc) {
			containingAnt = append(containingAnt, desc)
			continue
		}
		items = append(items, desc)
	}
	return append(containingAnt, items...)
}

func (dds deviceDescs) get(result string) *deviceDesc {
	for _, d := range dds {
		if result == d.humanReadable() {
			return d
		}
	}
	return nil
}

func getDeviceDescriptions(ctx *gousb.Context) (deviceDescs, error) {
	var ds []*deviceDesc
	_, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		ds = append(ds, (*deviceDesc)(desc))
		return false
	})
	return ds, errors.Wrap(err, "opening devices")
}

func userSelectDevice(ctx *gousb.Context) (*deviceDesc, error) {
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
