package usbprompt

import (
	"fmt"
	"regexp"

	"github.com/google/gousb"
	"github.com/google/gousb/usbid"
	"github.com/pkg/errors"
)

type DeviceDesc gousb.DeviceDesc

func (dd DeviceDesc) humanReadable() string {
	return fmt.Sprintf("%s (%s) (%s)",
		usbid.Vendors[dd.Vendor].Product[dd.Product],
		usbid.Vendors[dd.Vendor],
		usbid.Classes[dd.Class])
}

type deviceDescs []*DeviceDesc

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

func (dds deviceDescs) get(result string) *DeviceDesc {
	for _, d := range dds {
		if result == d.humanReadable() {
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
