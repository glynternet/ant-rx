package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/google/gousb"
	"github.com/google/gousb/usbid"
	"github.com/half2me/antgo/message"
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

func main() {
	if err := run(); err != nil {
		fmt.Printf("error whilst running: %+v", err)
	}
}

func run() error {
	ctx := gousb.NewContext()
	ctx.Debug(3)
	defer ctx.Close()

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
		return errors.Wrap(err, "running prompt")
	}

	fmt.Printf("You choose %q\n", result)
	chosen := ds.get(result)
	if chosen == nil {
		return errors.Errorf("chosen device cannot be found: %s", result)
	}

	dev, err := ctx.OpenDeviceWithVIDPID(chosen.Vendor, chosen.Product)
	if err != nil {
		return errors.Wrap(err, "opening device")
	}
	if dev == nil {
		return errors.New("open device is nil")
	}
	defer func() {
		if cErr := dev.Close(); cErr != nil {
			log.Printf("Error closing device: %v", cErr)
		}
	}()

	if err := dev.SetAutoDetach(true); err != nil {
		return errors.Wrap(err, "setting up autodetach")
	}

	cfg, err := dev.Config(1)
	if err != nil {
		return errors.Wrap(err, "something with configuration?!")
	}
	defer func() {
		if cErr := cfg.Close(); cErr != nil {
			log.Printf("Error closing config: %v", cErr)
		}
	}()

	intf, err := cfg.Interface(0, 0)
	if err != nil {
		return errors.Wrap(err, "claiming interface")
	}
	defer intf.Close()

	inep, err := intf.InEndpoint(1)
	if err != nil {
		return errors.Wrap(err, "using in-endpoint")
	}

	readstr, err := inep.NewStream(64, 1)
	if err != nil {
		return errors.Wrap(err, "preparing new in-endpoint stream for reading")
	}
	defer func() {
		if cErr := readstr.Close(); cErr != nil {
			log.Printf("Error closing in-endpoint stream: %v", cErr)
		}
	}()

	opCtx := context.Background()
	buf := make([]byte, 64)
	s := make(chan interface{})

	go func() {
		for {
			select {
			case <-s:
				return
			default:
				n, err := readstr.ReadContext(opCtx, buf)
				if err != nil {
					log.Printf("Error reading from stream: %+v", err)
					return
				}

				if n < 12 {
					continue
				}

				if buf[0] == message.MESSAGE_TX_SYNC {
					packet := message.AntPacket(buf)
					if packet.Class() == message.MESSAGE_TYPE_BROADCAST {
						msg := message.AntBroadcastMessage(packet)
						fmt.Println(msg.String())
					}
				}
			}
		}
	}()
	r := bufio.NewReader(os.Stdin)
	r.ReadLine()
	s <- true
	return nil
}
