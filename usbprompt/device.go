package usbprompt

import (
	"fmt"
	"regexp"

	"github.com/google/gousb"
	"github.com/google/gousb/usbid"
	"github.com/pkg/errors"
)

type DeviceDesc gousb.DeviceDesc

func (dd DeviceDesc) HumanReadable() string {
	return fmt.Sprintf("%s (%s) (%s)",
		usbid.Vendors[dd.Vendor].Product[dd.Product],
		usbid.Vendors[dd.Vendor],
		usbid.Classes[dd.Class])
}

type deviceDescs []*DeviceDesc

func (dds deviceDescs) matching(match func(*DeviceDesc) bool) deviceDescs {
	var matches deviceDescs
	for _, dd := range dds {
		if match(dd) {
			matches = append(matches, dd)
		}
	}
	return matches
}

func (dds deviceDescs) humanReadables() []string {
	var items []string
	for _, d := range dds {
		items = append(items, d.HumanReadable())
	}
	return items
}

func (dds deviceDescs) get(result string) *DeviceDesc {
	for _, d := range dds {
		if result == d.HumanReadable() {
			return d
		}
	}
	return nil
}

func getDeviceDescriptions(ctx *gousb.Context) (deviceDescs, error) {
	var ds []*DeviceDesc
	_, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		ds = append(ds, (*DeviceDesc)(desc))
		return false
	})
	return ds, errors.Wrap(err, "opening devices")
}

var _ = containsANTFunc()

func containsANTFunc() func(s string) bool {
	return regexp.MustCompile(`.*[aA][nN][tT].*`).MatchString
}
